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
	var _timetable, _email string
	var userID int
	timetable := fmt.Sprintf("%s,%s", fileID, sheetID)
	row := e.DB.QueryRow(`
		SELECT a.timetable_id, u.email, u.id
		FROM allowlist AS a, users AS u
		WHERE a.timetable_id=?
			AND a.user_id=u.id;
	`, timetable)
	if err := row.Scan(&_timetable, &_email, &userID); err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 404, "존재하지 않는 timetable")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 발생 : "+err.Error())
		return
	}
	if _email != email {
		functions.ResponseError(w, 403, "timetable 접근 권한 부족")
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
			functions.ResponseError(w, 500, err.Error())
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

	// Querying (Cell Validation Check)
	isPossible := true
	rows, err := e.DB.Query(`
		SELECT cell_start, cell_end
		FROM transactions
		WHERE transaction_type=1
			AND cell_column=?;
	`, *(reqData.Column))
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

	// Querying (Making a Transaction)
	res, err := e.DB.Exec(`
		INSERT INTO transactions (transaction_type, user_id, timetable_id, lecture, capacity, cell_column, cell_start, cell_end, professor)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?);
		`, userID, timetable, *(reqData.Lecture), *(reqData.Capacity), *(reqData.Column), *(reqData.Start), *(reqData.End), *(reqData.Professor))
	if err != nil {
		functions.ResponseError(w, 500, err.Error())
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

	// Check Timetable Permission
	var _timetable, _email string
	var userID int
	timetable := fmt.Sprintf("%s,%s", fileID, sheetID)
	row := e.DB.QueryRow(`
		SELECT a.timetable_id, u.email, u.id
		FROM allowlist AS a, users AS u
		WHERE a.timetable_id=?
			AND a.user_id=u.id;
	`, timetable)
	if err := row.Scan(&_timetable, &_email, &userID); err != nil {
		if err == sql.ErrNoRows {
			functions.ResponseError(w, 404, "존재하지 않는 timetable")
			return
		}
		functions.ResponseError(w, 500, "예기치 못한 에러 발생 : "+err.Error())
		return
	}
	if _email != email {
		functions.ResponseError(w, 403, "timetable 접근 권한 부족")
		return
	}

	// Check Transaction Permission
	var _transactionType int64
	row = e.DB.QueryRow(`
		SELECT u.email, t.transaction_type
		FROM transactions AS t, users AS u
		WHERE t.user_id=u.id
			AND t.transaction_id=?;
	`, reservationID)
	err := row.Scan(&_email, &_transactionType)
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
	res, err := e.DB.Exec(`
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

	functions.ResponseOK(w, "success", nil)
}
