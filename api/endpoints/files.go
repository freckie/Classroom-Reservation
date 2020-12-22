package endpoints

import (
	"bytes"
	"classroom/functions"
	"classroom/models"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GET /files
func (e *Endpoints) FilesGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get user email
	var email string
	if _email, ok := r.Header["X-User-Email"]; ok {
		email = _email[0]
	} else {
		functions.ResponseError(w, 401, "X-User-Email 헤더를 보내세요.")
		return
	}

	// Permission Check
	var isSuper int
	row := e.DB.QueryRow(`
		SELECT is_super FROM users WHERE email=?;
	`, email)
	if err := row.Scan(&isSuper); err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 401, "해당 유저가 존재하지 않습니다.")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	if isSuper == 0 {
		functions.ResponseError(w, 403, "접근 권한 부족. 관리자만 허용된 기능입니다.")
		return
	}

	// Result Resp
	resp := models.FilesGetResponse{}
	resp.Files = []models.FilesGetItem{}

	// Querying
	rows, err := e.DB.Query(`
		SELECT id, name, created_at FROM files ORDER BY created_at DESC;`)
	if err != nil {
		if err == sql.ErrNoRows {
			resp.FilesCount = 0
			functions.ResponseOK(w, "success", resp)
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var fileID, fileName, createdAtStr string
		err := rows.Scan(&fileID, &fileName, &createdAtStr)
		if err != nil {
			continue
		}

		temp := models.FilesGetItem{
			FileID:    fileID,
			FileName:  fileName,
			CreatedAt: functions.ToKST(createdAtStr),
		}
		resp.Files = append(resp.Files, temp)
	}

	// Struct for response
	resp.FilesCount = len(resp.Files)

	functions.ResponseOK(w, "success", resp)
}

// POST /files
func (e *Endpoints) FilesPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get user email
	var email string
	if _email, ok := r.Header["X-User-Email"]; ok {
		email = _email[0]
	} else {
		functions.ResponseError(w, 401, "X-User-Email 헤더를 보내세요.")
		return
	}

	// Permission Check
	var isSuper int
	row := e.DB.QueryRow(`
		SELECT is_super FROM users WHERE email=?;
	`, email)
	if err := row.Scan(&isSuper); err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 401, "해당 유저가 존재하지 않습니다.")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	if isSuper == 0 {
		functions.ResponseError(w, 403, "접근 권한 부족. 관리자만 허용된 기능입니다.")
		return
	}

	// Parse multipart/form-data
	r.ParseMultipartForm(10 << 20)

	// Parsing File
	var fileName string
	file, handler, err := r.FormFile("file")
	if err != nil {
		functions.ResponseError(w, 500, "파일 로드 중 예기치 못한 에러 : "+err.Error())
		return
	}
	defer file.Close()
	fileName = handler.Filename

	// Create temp file
	tempFile, err := ioutil.TempFile(os.TempDir(), "upload-*.tmp")
	if err != nil {
		functions.ResponseError(w, 500, "임시 파일 생성 중 예기치 못한 에러 : "+err.Error())
		return
	}
	defer tempFile.Close()

	// Write data to temp file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		functions.ResponseError(w, 500, "임시 파일 작성 중 예기치 못한 에러 : "+err.Error())
		return
	}
	tempFile.Write(fileBytes)

	// [Drive] Upload file to Google Drive
	reader := bytes.NewReader(fileBytes)
	sheetFile, err := e.Drive.UploadFile("KHU Classroom Reservation", fileName, reader)
	if err != nil {
		functions.ResponseError(w, 500, "구글 드라이브에 파일 업로드 중 예기치 못한 에러 : "+err.Error())
		return
	}

	// [Sheets] Get sheet properties of new file
	props, err := e.Sheets.GetAllSheetProperties(sheetFile.Id)
	if err != nil {
		functions.ResponseError(w, 500, "파일 속성 불러오기 중 예기치 못한 에러 : "+err.Error())
		return
	}

	// Querying with Transaction
	tx, err := e.DB.Begin()
	if err != nil {
		functions.ResponseError(w, 500, "트랜잭션 시작 중 예기치 못한 에러 : "+err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO files (id, name)
		VALUES (?, ?);
	`, sheetFile.Id, fileName)
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	query := "INSERT INTO sheets (id, name, file_id) VALUES "
	vals := []interface{}{}
	for idx := range props {
		query += "(?, ?, ?),"
		vals = append(vals, props[idx].SheetId, props[idx].Title, sheetFile.Id)
	}
	query = query[0 : len(query)-1]

	stmt, _ := tx.Prepare(query)
	_, err = stmt.Exec(vals...)
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	resp := models.FilesPostResponse{
		FileID: sheetFile.Id,
	}

	functions.ResponseOK(w, "success", resp)
}

// POST /files/<file_id>/share
func (e *Endpoints) FilesSharePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get user email
	var email string
	if _email, ok := r.Header["X-User-Email"]; ok {
		email = _email[0]
	} else {
		functions.ResponseError(w, 401, "X-User-Email 헤더를 보내세요.")
		return
	}

	// Get Path Parameters
	fileID := ps.ByName("file_id")

	// Permission Check
	var isSuper int
	row := e.DB.QueryRow(`
		SELECT is_super FROM users WHERE email=?;
	`, email)
	if err := row.Scan(&isSuper); err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 401, "해당 유저가 존재하지 않습니다.")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	if isSuper == 0 {
		functions.ResponseError(w, 403, "접근 권한 부족. 관리자만 허용된 기능입니다.")
		return
	}

	// Parse Request Data
	type reqDataStruct struct {
		UserEmails []string `json:"user_emails"`
	}
	var reqData reqDataStruct
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			functions.ResponseError(w, 500, err.Error())
		}
		json.Unmarshal(body, &reqData)
	} else {
		functions.ResponseError(w, 400, "JSON 형식만 가능합니다.")
		return
	}
	if len(reqData.UserEmails) == 0 {
		functions.ResponseError(w, 400, "user_emails를 한 개 이상 보내주세요.")
		return
	}

	// [Drive] Sharing file to users
	err := e.Drive.ShareFile(fileID, reqData.UserEmails)
	if err != nil {
		functions.ResponseError(w, 500, "파일 권한 설정 중 예기치 못한 에러 : "+err.Error())
		return
	}

	functions.ResponseOK(w, "success", nil)
}

// POST /files/<file_id>/protect
func (e *Endpoints) FilesProtectPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get user email
	var email string
	if _email, ok := r.Header["X-User-Email"]; ok {
		email = _email[0]
	} else {
		functions.ResponseError(w, 401, "X-User-Email 헤더를 보내세요.")
		return
	}

	// Get Path Parameters
	fileID := ps.ByName("file_id")

	// Permission Check
	var isSuper int
	row := e.DB.QueryRow(`
		SELECT is_super FROM users WHERE email=?;
	`, email)
	if err := row.Scan(&isSuper); err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 401, "해당 유저가 존재하지 않습니다.")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	if isSuper == 0 {
		functions.ResponseError(w, 403, "접근 권한 부족. 관리자만 허용된 기능입니다.")
		return
	}

	// [Sheets] Get sheet properties of new file
	props, err := e.Sheets.GetAllSheetProperties(fileID)
	if err != nil {
		functions.ResponseError(w, 500, "파일 속성 불러오기 중 예기치 못한 에러 : "+err.Error())
		return
	}

	var sheetIDs []int64
	for idx := range props {
		sheetIDs = append(sheetIDs, props[idx].SheetId)
	}

	// [Sheets] Protect all sheets
	err = e.Sheets.ProtectAll(fileID, sheetIDs)
	if err != nil {
		functions.ResponseError(w, 500, "셀 보호 설정 중 예기치 못한 에러 : "+err.Error())
		return
	}

	functions.ResponseOK(w, "success", nil)
}
