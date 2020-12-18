package endpoints

import (
	"classroom/functions"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// POST /timetables/<file_id>/<sheet_id>/allow
func (e *Endpoints) AllowlistPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		SELECT count(a.timetable_id), u.is_super
		FROM allowlist AS a, users AS u
		WHERE a.user_id=u.id
			AND a.timetable_id=?
			AND u.email=?;
	`, timetable, email)
	if err := row.Scan(&_count, &_isSuper); err == nil {
		if _count <= 0 {
			functions.ResponseError(w, 403, "timetable에 접근할 권한이 부족합니다.")
			return
		}
		if _isSuper != 1 {
			functions.ResponseError(w, 403, "관리자만 접근할 수 있는 기능입니다.")
			return
		}
	}

	// Parse Request Data
	type reqDataStruct struct {
		Email *string `json:"email"`
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
	if reqData.Email == nil {
		functions.ResponseError(w, 400, "파라미터를 전부 보내주세요.")
		return
	}

	// Querying
	_, err := e.DB.Exec(`
		INSERT INTO allowlist
		VALUES (?, (
			SELECT id FROM users WHERE email=?
		));
		`, timetable, *(reqData.Email))
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	functions.ResponseOK(w, "success", nil)
}
