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

// SheetsService is a wrapper for Spread Sheets Service and its context.
type SheetsService struct {
	srv *sheets.Service
	ctx context.Context
}

// NewSheetsService is a factory function which returns a new SheetsService{}.
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

// WriteAndMerge makes requests that merge cells and write value into cells.
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

// RemoveValue makes requests that clear  and unmerge cells.
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

func (s *SheetsService) GetAllSheetProperties(fileID string) ([]*sheets.SheetProperties, error) {
	var result []*sheets.SheetProperties

	req, err := s.srv.Spreadsheets.Get(fileID).Do()
	if err != nil {
		return nil, err
	}

	for idx := range req.Sheets {
		result = append(result, req.Sheets[idx].Properties)
	}

	return result, nil
}

func (s *SheetsService) ProtectAll(fileID string, sheetIDs []int64) error {
	reqs := []*sheets.Request{}

	for _, sheetID := range sheetIDs {
		req := &sheets.Request{}
		req.AddProtectedRange = &sheets.AddProtectedRangeRequest{
			ProtectedRange: &sheets.ProtectedRange{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartColumnIndex: 0,
					EndColumnIndex:   100,
					StartRowIndex:    0,
					EndRowIndex:      1000,
				},
			},
		}

		reqs = append(reqs, req)
	}

	rb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: reqs,
	}

	_, err := s.srv.Spreadsheets.BatchUpdate(fileID, rb).Context(s.ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

// SheetsRequest is a wrapper for request, especially WriteAndMerge() and RemoveValue() functions.
type SheetsRequest struct {
	SpreadSheetID string
	SheetName     string
	SheetID       int64
	Range         *sheets.GridRange
	RangeStr      string
	Value         string
}

// NewSheetsRequest is a factory function which returns a new SheetsRequest{}.
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
