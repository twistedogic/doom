package drive

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

const (
	FolderType = "application/vnd.google-apps.folder"
	TextType   = "text/plain"
)

func containId(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

type Drive struct {
	service  *drive.Service
	Base     string
	ParentID string
}

func New(credFile, cacheFile, base string) (*Drive, error) {
	ctx := context.TODO()
	client, err := GetClient(ctx, drive.DriveScope, credFile, cacheFile)
	if err != nil {
		return nil, err
	}
	srv, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	d := &Drive{
		service: srv,
		Base:    base,
	}
	id, err := d.MkdirAll(base)
	if err != nil {
		return nil, err
	}
	d.ParentID = id
	return d, nil
}

func (d *Drive) Find(name string) (*drive.FileList, error) {
	query := fmt.Sprintf("name='%s'", name)
	return d.service.Files.List().Q(query).Do()
}

func (d *Drive) List() ([]string, error) {
	ctx := context.TODO()
	listCall := d.service.Files.List()
	names := []string{}
	for {
		err := listCall.Pages(ctx, func(list *drive.FileList) error {
			for _, file := range list.Files {
				if containId(file.Parents, d.ParentID) {
					names = append(names, file.Name)
				}
			}
			return nil
		})
		if err != nil {
			return names, err
		}
	}
	return names, nil
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

func (d *Drive) Delete(name string) error {
	id, err := d.pathToId(name)
	if err != nil {
		return err
	}
	return d.service.Files.Delete(id).Do()
}

func (d *Drive) Download(id string, w io.Writer) error {
	fmt.Println(id)
	res, err := d.service.Files.Get(id).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(w, res.Body)
	return err
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

func (d *Drive) pathToId(path string) (string, error) {
	var fileId string
	list, err := d.Find(path)
	if err != nil {
		return fileId, err
	}
	for _, f := range list.Files {
		for _, p := range f.Parents {
			if p == d.ParentID {
				return f.Id, nil
			}
		}
	}
	return fileId, os.ErrNotExist
}

func (d *Drive) ReadFile(path string, w io.Writer) error {
	id, err := d.pathToId(path)
	if err != nil {
		return err
	}
	return d.Download(id, w)
}

func (d *Drive) WriteFile(path string, r io.ReadWriter) error {
	file := &drive.File{
		Name:     path,
		Parents:  []string{d.ParentID},
		MimeType: TextType,
	}
	id, err := d.pathToId(path)
	if os.IsNotExist(err) {
		_, err := d.Create(file, r)
		return err
	}
	if err != nil {
		return err
	}
	return d.Update(id, file, r)
}
