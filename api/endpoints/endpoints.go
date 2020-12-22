package endpoints

import (
	"classroom/utils"
	"database/sql"
)

type Endpoints struct {
	DB     *sql.DB
	Sheets *utils.SheetsService
	Drive  *utils.DriveService
}
