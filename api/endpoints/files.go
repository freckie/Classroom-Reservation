package endpoints

import (
	"classroom/functions"
	"classroom/models"
	"database/sql"
	"net/http"

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
