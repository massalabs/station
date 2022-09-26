// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// ThyraEventsGetterURL generates an URL for the thyra events getter operation
type ThyraEventsGetterURL struct {
	Caller string
	Str    string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ThyraEventsGetterURL) WithBasePath(bp string) *ThyraEventsGetterURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ThyraEventsGetterURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *ThyraEventsGetterURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/thyra/events/{str}/{caller}"

	caller := o.Caller
	if caller != "" {
		_path = strings.Replace(_path, "{caller}", caller, -1)
	} else {
		return nil, errors.New("caller is required on ThyraEventsGetterURL")
	}

	str := o.Str
	if str != "" {
		_path = strings.Replace(_path, "{str}", str, -1)
	} else {
		return nil, errors.New("str is required on ThyraEventsGetterURL")
	}

	_basePath := o._basePath
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *ThyraEventsGetterURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *ThyraEventsGetterURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *ThyraEventsGetterURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on ThyraEventsGetterURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on ThyraEventsGetterURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *ThyraEventsGetterURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
