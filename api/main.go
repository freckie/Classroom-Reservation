package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"classroom/endpoints"
	"classroom/functions"
	"classroom/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var logger *log.Logger

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		ip := r.RemoteAddr
		log.Println("[Req]", r.Method, r.URL, ip)
		handler.ServeHTTP(w, r)
	} else {
		functions.ResponseError(w, 403, "Forbidden hostname : "+r.Host)
	}
}

func main() {
	// Logger
	logger = log.New(os.Stdout, "LOG ", log.LstdFlags)

	// Config
	cfg, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// DB Setting
	if !(cfg.Server.LocalMode) {
		time.Sleep(time.Second * 10)
	}
	dbStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?allowNativePasswords=true&parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Schema)
	db, err := sql.Open("mysql", dbStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Google API Setting
	sheets, err := utils.NewSheetsService(cfg.Google.CredentialsPath)
	if err != nil {
		log.Fatal(err)
	}

	ep := endpoints.Endpoints{
		DB:     db,
		Sheets: sheets,
	}

	// Router Setting
	router := httprouter.New()
	router.GET("/api", ep.IndexGet)
	router.GET("/api/users", ep.UsersGet)
	router.POST("/api/users", ep.UsersPost)
	router.GET("/api/timetables/:file_id/:sheet_id/cell", ep.CellGet)
	router.POST("/api/timetables/:file_id/:sheet_id/allow", ep.AllowlistPost)
	router.DELETE("/api/timetables/:file_id/:sheet_id/allow", ep.AllowlistDelete)
	router.POST("/api/timetables/:file_id/:sheet_id/reservation", ep.ReservationPost)
	router.DELETE("/api/timetables/:file_id/:sheet_id/reservation/:reservation_id", ep.ReservationDelete)

	// Local Mode
	portStr := strconv.Itoa(cfg.Server.Port)
	if cfg.Server.LocalMode {
		handler := cors.AllowAll().Handler(router)
		hs := make(HostSwitch)
		hostname := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		hs[hostname] = handler
		hs["localhost:8000"] = handler

		// Start Server in Local Mode
		log.Println("[Local Mode] Starting HTTP API Server on port", portStr)
		log.Fatal(http.ListenAndServe(":"+portStr, hs))

	} else { // Release Mode
		handler := cors.AllowAll().Handler(router)
		hs := make(HostSwitch)
		hs["web-api"] = handler

		// Start Server
		log.Println("[Release Mode] Starting HTTP API Server on port", portStr)
		log.Fatal(http.ListenAndServe(":"+portStr, hs))
	}
}
