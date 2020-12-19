package endpoints

import (
	"classroom/functions"
	"classroom/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GET /files/<file_id>/<sheet_id>/cell
func (e *Endpoints) CellGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Get Query Parameters
	qp := r.URL.Query()
	var cellColumn string
	var cellStart int
	var cellEnd int
	var err error

	if _cellColumn, ok := qp["column"]; ok {
		cellColumn = strings.ToUpper(_cellColumn[0])
	} else {
		functions.ResponseError(w, 400, "column 파라미터를 보내세요.")
		return
	}

	if _cellStart, ok := qp["start"]; ok {
		cellStart, err = strconv.Atoi(_cellStart[0])
		if err != nil {
			functions.ResponseError(w, 400, "start 파라미터를 보내세요.")
			return
		}
	}
	if _cellEnd, ok := qp["end"]; ok {
		cellEnd, err = strconv.Atoi(_cellEnd[0])
		if err != nil {
			functions.ResponseError(w, 400, "end 파라미터를 보내세요.")
			return
		}
	}

	// Check Permission
	var _count int64
	var isSuper, sheetIDAuto sql.NullInt64
	row := e.DB.QueryRow(`
		SELECT
			(SELECT count(s.id)
			FROM sheets AS s, files AS f
			WHERE s.file_id=f.id
				AND f.id=?
				AND s.id=?) AS count,
			(SELECT id_auto
				FROM sheets
				WHERE id=?) AS id_auto,
			(SELECT is_super
			FROM users
			WHERE email=?) AS is_super;
	`, fileID, sheetID, sheetID, email)
	if err := row.Scan(&_count, &sheetIDAuto, &isSuper); err == nil {
		if _count != 1 || !sheetIDAuto.Valid {
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

	// Result Resp
	resp := models.CellGetResponse{}
	resp.Cells = []models.CellItem{}

	// Querying
	rows, err := e.DB.Query(`
		SELECT u.email, u.id, t.cell_column, t.cell_start, t.cell_end, t.lecture, t.professor, t.transaction_id, t.created_at, t.capacity
		FROM transactions AS t, users AS u
		WHERE t.user_id=u.id
			AND t.transaction_type=1
			AND t.sheet_id=?
			AND t.cell_column=?;`, sheetIDAuto.Int64, cellColumn)
	if err != nil {
		if err == sql.ErrNoRows {
			resp.CellsCount = 0
			functions.ResponseOK(w, "success", resp)
			return
		}
		functions.ResponseError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	cells := []models.CellTransactionModel{}
	for rows.Next() {
		temp := models.CellTransactionModel{}
		err := rows.Scan(&temp.UserEmail, &temp.UserID, &temp.CellColumn, &temp.CellStart, &temp.CellEnd, &temp.Lecture, &temp.Professor, &temp.TransactionID, &temp.CreatedAt, &temp.Capacity)
		if err != nil {
			continue
		}
		cells = append(cells, temp)
	}

	// Compare
	for i := cellStart; i <= cellEnd; i++ {
		isInRange := false
		for _, cell := range cells {
			if functions.InRange(i, cell.CellStart, cell.CellEnd) {
				temp := models.CellItem{}
				temp.Cell = fmt.Sprintf("%s%d", cellColumn, i)
				temp.IsReserved = true
				temp.UserEmail = cell.UserEmail
				temp.UserID = cell.UserID
				temp.Lecture = cell.Lecture
				temp.Professor = cell.Professor
				temp.TransactionID = cell.TransactionID
				temp.CreatedAt = functions.ToKST(cell.CreatedAt)
				temp.Capacity = cell.Capacity

				resp.Cells = append(resp.Cells, temp)
				isInRange = true
				break
			}
		}

		if !isInRange {
			temp := models.CellItem{}
			temp.Cell = fmt.Sprintf("%s%d", cellColumn, i)
			temp.IsReserved = false

			resp.Cells = append(resp.Cells, temp)
		}
	}

	// Struct for response
	resp.CellsCount = len(resp.Cells)

	functions.ResponseOK(w, "success", resp)
}
