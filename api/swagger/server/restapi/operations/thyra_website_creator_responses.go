// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// ThyraWebsiteCreatorOKCode is the HTTP code returned for type ThyraWebsiteCreatorOK
const ThyraWebsiteCreatorOKCode int = 200

/*ThyraWebsiteCreatorOK Page found

swagger:response thyraWebsiteCreatorOK
*/
type ThyraWebsiteCreatorOK struct {
}

// NewThyraWebsiteCreatorOK creates ThyraWebsiteCreatorOK with default headers values
func NewThyraWebsiteCreatorOK() *ThyraWebsiteCreatorOK {

	return &ThyraWebsiteCreatorOK{}
}

// WriteResponse to the client
func (o *ThyraWebsiteCreatorOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}
