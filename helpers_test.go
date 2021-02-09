package gouploader

import (
	"fmt"
	"github.com/golang/mock/gomock"
	m "github.com/paulhenri-l/gouploader/mocks/contracts"
	"os"
	"testing"
)

func fakeFile(t *testing.T, contents string) string {
	tmp := t.TempDir()
	fp := fmt.Sprintf("%s/%s", tmp, "fake_file")

	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(contents))
	if err != nil {
		panic(err)
	}

	return fp
}

func fakeUploader(t *testing.T) (*m.MockUploader, *gomock.Controller) {
	ctl := gomock.NewController(t)
	t.Cleanup(func() {
		ctl.Finish()
	})

	return m.NewMockUploader(ctl), ctl
}
