package endpoints

import (
	"classroom/functions"
	"classroom/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// POST /timetables/<file_id>/<sheet_id>/reservation
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
	var _count, _isSuper int64
	timetable := fmt.Sprintf("%s,%s", fileID, sheetID)
	row := e.DB.QueryRow(`
		SELECT count(timetable_id)
		FROM allowlist
		WHERE timetable_id=?;
	`, timetable)
	if err := row.Scan(&_count); err == nil {
		if _count <= 0 {
			functions.ResponseError(w, 404, "존재하지 않는 timetable.")
			return
		}
	}

	row = e.DB.QueryRow(`
		SELECT (
			SELECT count(a.timetable_id)
			FROM allowlist AS a, users AS u
			WHERE a.user_id=u.id
				AND a.timetable_id=?
				AND u.email=?
		) AS count,
		(
			SELECT is_super FROM users WHERE email=?
		) AS is_super;
	`, timetable, email, email)
	if err := row.Scan(&_count, &_isSuper); err == nil {
		if _isSuper != 1 {
			functions.ResponseError(w, 403, "관리자만 접근할 수 있는 기능입니다.")
			return
		}
		if _count <= 0 {
			functions.ResponseError(w, 403, "timetable에 접근할 권한이 부족합니다.")
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
			AND timetable_id=?
			AND cell_column=?;
	`, timetable, *(reqData.Column))
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
		INSERT INTO transactions (transaction_type, user_id, timetable_id, lecture, capacity, cell_column, cell_start, cell_end, professor)
		VALUES (1, (
			SELECT id FROM users WHERE email=?
		), ?, ?, ?, ?, ?, ?, ?);
		`, email, timetable, *(reqData.Lecture), *(reqData.Capacity), *(reqData.Column), *(reqData.Start), *(reqData.End), *(reqData.Professor))
	if err != nil {
		functions.ResponseError(w, 500, err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	resp.TransactionID, err = res.LastInsertId()
	if err != nil {
		functions.ResponseError(w, 500, err.Error())
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

// DELETE /timetables/<file_id>/<sheet_id>/reservation/<reservation_id>
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
	var _count, _isSuper int64
	timetable := fmt.Sprintf("%s,%s", fileID, sheetID)
	row := e.DB.QueryRow(`
		SELECT count(timetable_id)
		FROM allowlist
		WHERE timetable_id=?;
	`, timetable)
	if err := row.Scan(&_count); err == nil {
		if _count <= 0 {
			functions.ResponseError(w, 404, "존재하지 않는 timetable.")
			return
		}
	}

	row = e.DB.QueryRow(`
		SELECT (
			SELECT count(a.timetable_id)
			FROM allowlist AS a, users AS u
			WHERE a.user_id=u.id
				AND a.timetable_id=?
				AND u.email=?
		) AS count,
		(
			SELECT is_super FROM users WHERE email=?
		) AS is_super;
	`, timetable, email, email)
	if err := row.Scan(&_count, &_isSuper); err == nil {
		if _isSuper != 1 {
			functions.ResponseError(w, 403, "관리자만 접근할 수 있는 기능입니다.")
			return
		}
		if _count <= 0 {
			functions.ResponseError(w, 403, "timetable에 접근할 권한이 부족합니다.")
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
	var _transactionType int64
	var _email string
	row = tx.QueryRow(`
		SELECT u.email, t.transaction_type
		FROM transactions AS t, users AS u
		WHERE t.user_id=u.id
			AND t.transaction_id=?;
	`, reservationID)
	err = row.Scan(&_email, &_transactionType)
	if err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 404, "존재하지 않는 예약")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 발생 : "+err.Error())
		return
	}
	if _email != email {
		functions.ResponseError(w, 403, "예약 접근 권한 부족")
		return
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

	err = tx.Commit()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	functions.ResponseOK(w, "success", nil)
}
