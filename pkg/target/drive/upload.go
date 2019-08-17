package drive

import (
	"os"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type Drive struct {
	*drive.Service
}

func NewDrive(credFile, cacheFile string) (*Drive, error) {
	ctx := context.Background()
	scope := drive.DriveScope
	client, err := GetClient(ctx, scope, credFile, cacheFile)
	if err != nil {
		return nil, err
	}
	srv, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	return &Drive{srv}, nil
}

func (d *Drive) UploadCSV(filename string) error {
	baseMimeType := "text/csv"
	convertedMimeType := "application/vnd.google-apps.spreadsheet"
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	f := &drive.File{
		Name:     filename,
		MimeType: convertedMimeType,
	}
	createFile := d.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType))
	if _, err := createFile.Do(); err != nil {
		return err
	}
	return nil
}
