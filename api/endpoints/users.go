package endpoints

import (
	"classroom/functions"
	"classroom/models"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GET /users
func (e *Endpoints) UsersGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			functions.ResponseError(w, 401, "해당 유저가 존재하지 않음")
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
	resp := models.UsersGetResponse{}
	resp.Users = []models.UsersGetItem{}

	// Querying
	rows, err := e.DB.Query(`
		SELECT id, email, is_super FROM users;`)
	if err != nil {
		if err == sql.ErrNoRows {
			resp.UsersCount = 0
			functions.ResponseOK(w, "success", resp)
			return
		}
		functions.ResponseError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		temp := models.UsersGetItem{}
		err := rows.Scan(&temp.UserID, &temp.UserEmail, &temp.IsSuper)
		if err != nil {
			continue
		}
		resp.Users = append(resp.Users, temp)
	}

	// Struct for response
	resp.UsersCount = len(resp.Users)

	functions.ResponseOK(w, "success", resp)
}

// POST /users
func (e *Endpoints) UsersPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get user email
	var email string
	if _email, ok := r.Header["X-User-Email"]; ok {
		email = _email[0]
	} else {
		functions.ResponseError(w, 401, "X-User-Email 헤더를 보내세요.")
		return
	}

	// Check Permission
	var _count int64
	row := e.DB.QueryRow(`
		SELECT count(id)
		FROM users
		WHERE is_super=1
			AND email=?;
	`, email)
	if err := row.Scan(&_count); err == nil {
		if _count <= 0 {
			functions.ResponseError(w, 403, "관리자 권한 부족.")
			return
		}
	} else {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	// Parse Request Data
	var isSuper bool
	type reqDataStruct struct {
		Email   *string `json:"email"`
		IsSuper *bool   `json:"is_super"`
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
	if reqData.Email == nil {
		functions.ResponseError(w, 400, "파라미터를 전부 보내주세요.")
		return
	}

	if reqData.IsSuper == nil {
		isSuper = false
	} else {
		isSuper = *(reqData.IsSuper)
	}

	// Querying with Transaction
	tx, err := e.DB.Begin()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}
	defer tx.Rollback()

	row = tx.QueryRow(`
		SELECT count(id)
		FROM users
		WHERE email=?;
	`, *(reqData.Email))
	if err = row.Scan(&_count); err == nil {
		if _count > 0 {
			functions.ResponseError(w, 500, "이미 존재하는 이메일입니다.")
			return
		}
	}

	result, err := tx.Exec(`
		INSERT INTO users (email, is_super)
		VALUES (?, ?);
	`, *(reqData.Email), isSuper)
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	// Result
	resp := models.UsersPostResponse{}
	resp.UserID, err = result.LastInsertId()
	if err != nil {
		functions.ResponseError(w, 500, "예기치 못한 에러 : "+err.Error())
		return
	}

	functions.ResponseOK(w, "success", resp)
	return
}
