package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// DriveService is a wrapper for Google Drive Service and its context.
type DriveService struct {
	srv               *drive.Service
	ctx               context.Context
	driveRootFolderID string
}

// NewDriveService is a factory function which returns a new DriveService{}.
func NewDriveService(credentialsPath, driveRootFolderID string) (*DriveService, error) {
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b,
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/spreadsheets",
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	ctx := context.Background()
	driveService, err := drive.NewService(ctx,
		option.WithHTTPClient(client),
		option.WithScopes(drive.DriveScope),
	)
	if err != nil {
		return nil, err
	}

	srv := &DriveService{
		srv:               driveService,
		ctx:               ctx,
		driveRootFolderID: driveRootFolderID,
	}

	return srv, nil
}

// UploadFile makes a request that uploads file to specific folder.
func (s *DriveService) UploadFile(folderName, fileName string, content io.Reader) (*drive.File, error) {
	file, err := s.createFile(fileName, s.driveRootFolderID, content)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ShareFile makes a request that shares the file for specific users.
func (s *DriveService) ShareFile(fileID string, emails []string) error {
	for _, email := range emails {
		perm := &drive.Permission{
			Type:         "user",
			Role:         "writer",
			EmailAddress: email,
		}
		resp, err := s.srv.Permissions.Create(fileID, perm).Do()
		if err != nil {
			fmt.Println(resp, err)
			continue
		}
	}

	return nil
}

// createFolder creates a folder and returns its object.
func (s *DriveService) createFolder(name, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}

	file, err := s.srv.Files.Create(d).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

// createFile creates a file and returns its object.
func (s *DriveService) createFile(name, parentID string, content io.Reader) (*drive.File, error) {
	f := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.spreadsheet",
		Parents:  []string{parentID},
	}

	file, err := s.srv.Files.Create(f).Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

type UploadRequest struct {
	FolderName string
	FileName   string
	MimeType   string
	Content    io.Reader
}
