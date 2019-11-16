package drive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

const (
	FolderType = "application/vnd.google-apps.folder"
)

type Drive struct {
	service    *drive.Service
	MimeType   string
	Credential string
	CacheFile  string
	Append     bool
}

func New(credFile, cacheFile string) (*Drive, error) {
	d := &Drive{
		Credential: credFile,
		CacheFile:  cacheFile,
	}
	err := d.Load()
	return d, err
}

func (d *Drive) Load() error {
	ctx := context.Background()
	client, err := GetClient(ctx, drive.DriveScope, d.Credential, d.CacheFile)
	if err != nil {
		return err
	}
	srv, err := drive.New(client)
	if err != nil {
		return err
	}
	d.service = srv
	return nil
}

func (d *Drive) Find(name string) (*drive.FileList, error) {
	query := fmt.Sprintf("name='%s'", name)
	return d.service.Files.List().Q(query).Do()
}

func (d *Drive) IsExist(name string) (id string, err error) {
	list, err := d.Find(name)
	if err != nil {
		return
	}
	switch {
	case len(list.Files) == 1:
		id = list.Files[0].Id
	case len(list.Files) == 0:
		err = os.ErrNotExist
	default:
		err = fmt.Errorf("%d files return for %s", len(list.Files), name)
	}
	return
}

func (d *Drive) Delete(id string) error {
	return d.service.Files.Delete(id).Do()
}

func (d *Drive) Download(id string, w io.Writer) error {
	res, err := d.service.Files.Get(id).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}
	return nil
}

func (d *Drive) Update(id string, file *drive.File, r io.Reader) error {
	_, err := d.service.Files.Update(id, file).Media(r).Do()
	return err
}

func (d *Drive) Create(file *drive.File, r io.Reader) (*drive.File, error) {
	return d.service.Files.Create(file).Media(r).Do()
}

func (d *Drive) MkdirAll(path string) (string, error) {
	var parentID string
	folders := strings.Split(path, string(os.PathSeparator))
	for _, name := range folders {
		folder := &drive.File{Name: name, MimeType: FolderType}
		id, err := d.IsExist(name)
		switch {
		case os.IsNotExist(err):
			if parentID != "" {
				folder.Parents = []string{parentID}
			}
			file, err := d.service.Files.Create(folder).Do()
			if err != nil {
				return parentID, err
			}
			parentID = file.Id
		case err != nil:
			return parentID, err
		default:
			parentID = id
		}
	}
	return parentID, nil
}

func (d *Drive) WriteFile(path string, r io.ReadWriter) (string, error) {
	var id string
	dirs, filename := filepath.Split(path)
	file := &drive.File{Name: filename, MimeType: d.MimeType}
	parentID, err := d.MkdirAll(dirs)
	if err != nil {
		return id, err
	}
	list, err := d.Find(filename)
	if err != nil {
		return id, err
	}
	for _, f := range list.Files {
		for _, p := range f.Parents {
			if p == parentID {
				if d.Append {
					if err := d.Download(f.Id, r); err != nil {
						return f.Id, err
					}
				}
				err := d.Update(f.Id, file, r)
				return f.Id, err
			}
		}
	}
	if parentID != "" {
		file.Parents = []string{parentID}
	}
	f, err := d.Create(file, r)
	if err != nil {
		return id, err
	}
	return f.Id, nil
}
