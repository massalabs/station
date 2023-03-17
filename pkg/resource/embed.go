package resource

import (
	"embed"
	"io/fs"
)

type embeddedResources struct {
	fs           embed.FS
	subDirectory string
}

func (r *embeddedResources) Read(file string) ([]byte, error) {
	content, err := r.fs.ReadFile(r.subDirectory + "/" + file)
	if err == fs.ErrNotExist {
		return nil, ErrNotExist
	}

	if err != nil {
		return nil, err
	}

	return content, nil
}

func NewEmbedStorer(fs embed.FS, subDirectory string) *embeddedResources {
	return &embeddedResources{fs: fs, subDirectory: subDirectory}
}

var _ Reader = new(embeddedResources)
