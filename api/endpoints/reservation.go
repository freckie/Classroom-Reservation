package endpoints

import (
	"classroom/functions"
	"classroom/models"
	"classroom/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// POST /files/<file_id>/<sheet_id>/reservation
func (e *Endpoints) ReservationPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	sheetID := ps.ByName("sheet_id")

	// Check Permission
	var _count int64
	var isSuper, sheetIDAuto sql.NullInt64
	var sheetName sql.NullString
	row := e.DB.QueryRow(`
		SELECT
			(SELECT count(s.id)
			FROM sheets AS s, files AS f
			WHERE s.file_id=f.id
				AND f.id=?
				AND s.id=?) AS count,
			(SELECT name
				FROM sheets
				WHERE id=?) AS sheet_name,
			(SELECT id_auto
				FROM sheets
				WHERE id=?) AS id_auto,
			(SELECT is_super
			FROM users
			WHERE email=?) AS is_super;
	`, fileID, sheetID, sheetID, sheetID, email)
	if err := row.Scan(&_count, &sheetName, &sheetIDAuto, &isSuper); err == nil {
		if _count != 1 || !sheetName.Valid || !sheetIDAuto.Valid {
			functions.ResponseError(w, 404, "해당 파일이나 시트가 존재하지 않습니다.")
			return
		}
		if !isSuper.Valid {
			functions.ResponseError(w, 401, "등록되지 않은 사용자입니다.")
			return
		}
	} else {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	// Parse Request Data
	type reqDataStruct struct {
		Column    *string `json:"column"`
		Start     *int    `json:"start"`
		End       *int    `json:"end"`
		Lecture   *string `json:"lecture"`
		Professor *string `json:"professor"`
		Capacity  *int    `json:"capacity"`
	}
	var reqData reqDataStruct
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
			return
		}
		json.Unmarshal(body, &reqData)
	} else {
		functions.ResponseError(w, 400, "JSON 형식만 가능합니다.")
		return
	}
	if reqData.Column == nil || reqData.Start == nil || reqData.End == nil ||
		reqData.Lecture == nil || reqData.Professor == nil || reqData.Capacity == nil {
		functions.ResponseError(w, 400, "파라미터를 전부 보내주세요.")
		return
	}

	// Querying with Transaction
	tx, err := e.DB.Begin()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	defer tx.Rollback()

	// Querying (Cell Validation Check)
	isPossible := true
	rows, err := tx.Query(`
		SELECT cell_start, cell_end
		FROM transactions
		WHERE transaction_type=1
			AND sheet_id=?
			AND cell_column=?;
	`, sheetIDAuto.Int64, *(reqData.Column))
	if err == sql.ErrNoRows {
		isPossible = true
	}
	defer rows.Close()

loopCheckingValidation:
	for rows.Next() {
		var _start, _end int
		err = rows.Scan(&_start, &_end)
		if err != nil {
			continue
		}
		for i := *(reqData.Start); i <= *(reqData.End); i++ {
			if functions.InRange(i, _start, _end) {
				isPossible = false
				break loopCheckingValidation
			}
		}
	}

	if !isPossible {
		functions.ResponseError(w, 500, "해당 셀 범위에 예약이 존재합니다.")
		return
	}

	// Result Resp
	resp := models.ReservationPostResponse{}

	// Querying
	res, err := tx.Exec(`
		INSERT INTO transactions (transaction_type, user_id, sheet_id, lecture, capacity, cell_column, cell_start, cell_end, professor)
		VALUES (1, (
			SELECT id FROM users WHERE email=?
		), ?, ?, ?, ?, ?, ?, ?);
		`, email, sheetIDAuto.Int64, *(reqData.Lecture), *(reqData.Capacity), *(reqData.Column), *(reqData.Start), *(reqData.End), *(reqData.Professor))
	if err != nil {
		functions.ResponseError(w, 500, err.Error())
		return
	}

	// Merge and write value on cells
	cellValue := fmt.Sprintf("%s\n%s", *(reqData.Lecture), *(reqData.Professor))
	sheetIDint, _ := strconv.Atoi(sheetID)
	sr := utils.NewSheetsRequest(
		fileID,
		sheetName.String,
		int64(sheetIDint),
		*(reqData.Column),
		int64(*(reqData.Start)),
		int64(*(reqData.End)),
		cellValue,
	)
	err = e.Sheets.WriteAndMerge(sr)
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	// Transaction Commit
	err = tx.Commit()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	resp.TransactionID, err = res.LastInsertId()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	resp.IsSuccess = true
	resp.CellColumn = *(reqData.Column)
	resp.CellStart = *(reqData.Start)
	resp.CellEnd = *(reqData.End)
	resp.Lecture = *(reqData.Lecture)
	resp.Professor = *(reqData.Professor)

	functions.ResponseOK(w, "success", resp)
}

// DELETE /files/<file_id>/<sheet_id>/reservation/<reservation_id>
func (e *Endpoints) ReservationDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	sheetID := ps.ByName("sheet_id")
	reservationID := ps.ByName("reservation_id")

	// Check Permission
	var _count int64
	var isSuper sql.NullInt64
	var sheetName sql.NullString
	row := e.DB.QueryRow(`
		SELECT
			(SELECT count(s.id)
			FROM sheets AS s, files AS f
			WHERE s.file_id=f.id
				AND f.id=?
				AND s.id=?) AS count,
			(SELECT name
				FROM sheets
				WHERE id=?) AS sheet_name,
			(SELECT is_super
			FROM users
			WHERE email=?) AS is_super;
	`, fileID, sheetID, sheetID, email)
	if err := row.Scan(&_count, &sheetName, &isSuper); err == nil {
		if _count != 1 || !sheetName.Valid {
			functions.ResponseError(w, 404, "해당 파일이나 시트가 존재하지 않습니다.")
			return
		}
		if !isSuper.Valid {
			functions.ResponseError(w, 401, "등록되지 않은 사용자입니다.")
			return
		}
	} else {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	// Querying with Transaction
	tx, err := e.DB.Begin()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	defer tx.Rollback()

	// Check Transaction Permission
	var _transactionType, _cellStart, _cellEnd int64
	var _email, _cellColumn string
	row = tx.QueryRow(`
		SELECT u.email, t.transaction_type, t.cell_column, t.cell_start, t.cell_end
		FROM transactions AS t, users AS u
		WHERE t.user_id=u.id
			AND t.transaction_id=?;
	`, reservationID)
	err = row.Scan(&_email, &_transactionType, &_cellColumn, &_cellStart, &_cellEnd)
	if err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 404, "존재하지 않는 예약")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 발생 : "+err.Error())
		return
	}
	if isSuper.Int64 == 0 {
		if _email != email {
			functions.ResponseError(w, 403, "예약 접근 권한 부족")
			return
		}
	}
	if _transactionType == 0 {
		functions.ResponseError(w, 500, "이미 취소된 예약")
		return
	}

	// Querying
	res, err := tx.Exec(`
		UPDATE transactions SET transaction_type=0 WHERE transaction_id=?
	`, reservationID)
	if err != nil {
		functions.ResponseError(w, 500, err.Error())
		return
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		functions.ResponseError(w, 500, "예기치 못한 에러 발생. (RowsAffected != 1)")
		return
	}

	// Unmerge and clear value on cells
	sheetIDint, _ := strconv.Atoi(sheetID)
	sr := utils.NewSheetsRequest(
		fileID,
		sheetName.String,
		int64(sheetIDint),
		_cellColumn,
		_cellStart,
		_cellEnd,
		"",
	)
	err = e.Sheets.RemoveValue(sr)
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	// Transaction Commit
	err = tx.Commit()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	functions.ResponseOK(w, "success", nil)
}
