package gouploader

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"path"
	"testing"
)

func TestGcs_Upload(t *testing.T) {
	c, _ := storage.NewClient(context.Background())
	b := c.Bucket("paulhenri-l-testing")
	f := fakeFile(t, "hello")
	uploader := NewGcs(b)

	res := uploader.Upload(f)
	assert.NoError(t, res.GetError())

	objs := getTestBucketObjects(t)
	assert.Equal(t, 1, len(objs))
	assert.Equal(t, path.Base(f), objs[0].Name)

	cleanTestBucket(t)
}

func TestGcs_Upload_BadFile(t *testing.T) {
	c, _ := storage.NewClient(context.Background())
	b := c.Bucket("paulhenri-l-testing")
	uploader := NewGcs(b)

	res := uploader.Upload("i-dont-exists")

	assert.Error(t, res.GetError())
	assert.Contains(t, res.GetError().Error(), "cannot get file stat")
}

func TestGcs_Upload_GcsError(t *testing.T) {
	c, _ := storage.NewClient(context.Background())
	b := c.Bucket("this-is-not-my-bucket")
	f := fakeFile(t, "hola")
	uploader := NewGcs(b)

	res := uploader.Upload(f)

	assert.Error(t, res.GetError())
	assert.Contains(t, res.GetError().Error(), "Gcs error")
}

func getTestBucketObjects(t *testing.T) []*storage.ObjectAttrs {
	var objs []*storage.ObjectAttrs
	c, _ := storage.NewClient(context.Background())
	b :=  c.Bucket("paulhenri-l-testing")
	defer c.Close()

	it := b.Objects(context.Background(), &storage.Query{})
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			t.Error(err)
		}

		objs = append(objs, objAttrs)
	}

	return objs
}

func cleanTestBucket(t *testing.T) {
	c, _ := storage.NewClient(context.Background())
	b :=  c.Bucket("paulhenri-l-testing")
	defer c.Close()

	objs := getTestBucketObjects(t)

	for _, obj := range objs {
		err := b.Object(obj.Name).Delete(context.Background())
		if err != nil {
			t.Error(err)
		}
	}
}
