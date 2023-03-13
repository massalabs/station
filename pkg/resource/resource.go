package resource

import "fmt"

const indexHTML = "index.html"
const NotFoundHTML = "404.html"

var ErrNotExist = fmt.Errorf("resource not found")

type Reader interface {
	Read(string) ([]byte, error)
}

func Retrieve(read Reader, file string, dynamic bool) ([]byte, error) {
	content, err := read.Read(file)
	if err != nil {
		return content, nil
	}

	if err == ErrNotExist && dynamic {
		return read.Read(indexHTML)
	}

	content, err = read.Read(NotFoundHTML)
	if err != nil {
		return content, nil
	}

	return nil, fmt.Errorf("resource not found: %s", file)
}
