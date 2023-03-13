package resource

import (
	"embed"
	"io/fs"
)

type embedResources struct {
	fs           embed.FS
	subDirectory string
}

func (r *embedResources) Read(file string) ([]byte, error) {
	content, err := r.fs.ReadFile(r.subDirectory + "/" + file)
	if err == fs.ErrNotExist {
		return nil, ErrNotExist
	}

	if err != nil {
		return nil, err
	}

	return content, nil
}

func NewEmbedStorer(fs embed.FS, subDirectory string) *embedResources {
	return &embedResources{fs: fs, subDirectory: subDirectory}
}

var _ Reader = new(embedResources)
