package gouploader

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/paulhenri-l/gouploader/entities"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"time"
)

type Gcs struct {
	bucket *storage.BucketHandle
}

func NewGcs(b *storage.BucketHandle) *Gcs {
	return &Gcs{bucket: b}
}

func (u *Gcs) Upload(file string) *entities.UploadResult {
	ur := &entities.UploadResult{Filepath: file}

	fi, err := os.Stat(file)
	if err != nil {
		ur.Error = errors.Wrap(err, "cannot get file stat")
		return ur
	}

	ur.Size = fi.Size()

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		ur.Error = errors.Wrap(err, "cannot open file")
		return ur
	}

	ur.Start = time.Now()
	err = u.doUpload(f, path.Base(file))
	ur.End = time.Now()
	if err != nil {
		ur.Error = errors.Wrap(err, "Gcs error")
		return ur
	}

	return ur
}

func (g *Gcs) doUpload(f io.Reader, object string) error {
	var err error

	wc := g.bucket.Object(object).NewWriter(context.Background())
	if _, err = io.Copy(wc, f); err != nil {
		return errors.Wrap(err, "error while copying file")
	}

	if err = wc.Close(); err != nil {
		return errors.Wrap(err, "error while closing")
	}

	return nil
}
