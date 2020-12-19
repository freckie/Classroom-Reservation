package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	srv *sheets.Service
	ctx context.Context
}

func NewSheetsService(credentialsPath string) (*SheetsService, error) {
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx,
		option.WithHTTPClient(client),
		option.WithScopes(sheets.SpreadsheetsScope),
	)
	if err != nil {
		return nil, err
	}

	srv := &SheetsService{
		srv: sheetsService,
		ctx: ctx,
	}

	return srv, nil
}

func (s *SheetsService) WriteAndMerge(sr SheetsRequest) error {
	req := &sheets.Request{}
	req.MergeCells = &sheets.MergeCellsRequest{
		MergeType: "MERGE_ALL",
		Range:     sr.Range,
	}

	rb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{req},
	}

	_, err := s.srv.Spreadsheets.BatchUpdate(sr.SpreadSheetID, rb).Context(s.ctx).Do()
	if err != nil {
		return err
	}

	var vr sheets.ValueRange
	val := []interface{}{sr.Value}
	vr.Values = append(vr.Values, val)
	_, err = s.srv.Spreadsheets.Values.Append(sr.SpreadSheetID, sr.RangeStr, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SheetsService) RemoveValue(sr SheetsRequest) error {
	req := &sheets.Request{}
	req.UnmergeCells = &sheets.UnmergeCellsRequest{
		Range: sr.Range,
	}

	req2 := &sheets.Request{}
	req2.UpdateBorders = &sheets.UpdateBordersRequest{
		Range: sr.Range,
		InnerHorizontal: &sheets.Border{
			Color: &sheets.Color{
				Blue:  0.0,
				Green: 0.0,
				Red:   0.0,
			},
			Style: "SOLID",
		},
	}

	rb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{req, req2},
	}

	_, err := s.srv.Spreadsheets.Values.Clear(sr.SpreadSheetID, sr.RangeStr, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return err
	}

	_, err = s.srv.Spreadsheets.BatchUpdate(sr.SpreadSheetID, rb).Context(s.ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

type SheetsRequest struct {
	SpreadSheetID string
	SheetName     string
	SheetID       int64
	Range         *sheets.GridRange
	RangeStr      string
	Value         string
}

func NewSheetsRequest(
	spreadSheetID, sheetName string, sheetID int64, column string, start, end int64, value string) SheetsRequest {
	colIndex := A1ToInt(column)

	req := SheetsRequest{}
	req.SpreadSheetID = spreadSheetID
	req.SheetID = sheetID
	req.Range = &sheets.GridRange{
		SheetId:          sheetID,
		StartColumnIndex: colIndex,
		EndColumnIndex:   colIndex + 1,
		StartRowIndex:    start - 1,
		EndRowIndex:      end,
	}
	req.RangeStr = fmt.Sprintf("%s!%s%d:%s%d", sheetName, column, start, column, end)
	req.Value = value
	return req
}
