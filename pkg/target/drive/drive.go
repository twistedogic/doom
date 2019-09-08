package drive

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/helper/token"
	"github.com/twistedogic/doom/pkg/target/ndjson"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

const (
	FolderType = "application/vnd.google-apps.folder"
)

type Drive struct {
	service    *drive.Service
	buf        *ndjson.Target
	Filename   string
	Credential string
	CacheFile  string
}

func New(filename, credFile, cacheFile string) (*Drive, error) {
	d := &Drive{Filename: filename, Credential: credFile, CacheFile: cacheFile}
	if err := d.load(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Drive) load() error {
	ctx := context.Background()
	buf := &ndjson.Target{File: &bytes.Buffer{}}
	client, err := token.GetClient(ctx, drive.DriveScope, d.Credential, d.CacheFile)
	if err != nil {
		return err
	}
	srv, err := drive.New(client)
	if err != nil {
		return err
	}
	d.buf = buf
	d.service = srv
	return nil
}

func (d *Drive) Load(c config.Setting) error {
	if err := c.ParseConfig(d); err != nil {
		return err
	}
	return d.load()
}

func (d *Drive) UpsertItem(i interface{}) error {
	return d.buf.UpsertItem(i)
}

func (d *Drive) Find(name string) (*drive.FileList, error) {
	query := fmt.Sprintf("name='%s'", name)
	return d.service.Files.List().Q(query).Do()
}

func (d *Drive) IsExist(name string) (id string, exist bool, err error) {
	list, err := d.Find(name)
	if err != nil {
		return
	}
	switch {
	case len(list.Files) == 1:
		id = list.Files[0].Id
		exist = true
	case len(list.Files) == 0:
		exist = false
	default:
		err = fmt.Errorf("%d files return for %s", len(list.Files), name)
	}
	return
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

func (d *Drive) createFolders(names ...string) (string, error) {
	var parentID string
	for _, name := range names {
		folder := &drive.File{
			Name:     name,
			MimeType: FolderType,
		}
		id, exist, err := d.IsExist(name)
		switch {
		case err != nil:
			return parentID, err
		case exist:
			parentID = id
		default:
			if parentID != "" {
				folder.Parents = []string{parentID}
			}
			file, err := d.service.Files.Create(folder).Do()
			if err != nil {
				return parentID, err
			}
			parentID = file.Id
		}
	}
	return parentID, nil
}

func (d *Drive) writeFile(path string) error {
	tokens := strings.Split(path, "/")
	dirs, filename := tokens[:len(tokens)-1], tokens[len(tokens)-1]
	file := &drive.File{Name: filename}
	parentID, err := d.createFolders(dirs...)
	if err != nil {
		return err
	}
	list, err := d.Find(filename)
	if err != nil {
		return err
	}
	for _, f := range list.Files {
		for _, p := range f.Parents {
			if p == parentID {
				update := d.service.Files.Update(f.Id, file).Media(d.buf.File)
				_, err := update.Do()
				return err
			}
		}
	}
	_, err = d.service.Files.Create(file).Media(d.buf.File).Do()
	return err
}

func (d *Drive) Delete(id string) error {
	return d.service.Files.Delete(id).Do()
}

func (d *Drive) WriteFile() error {
	return d.writeFile(d.Filename)
}

func (d *Drive) Close() error {
	return d.WriteFile()
}
