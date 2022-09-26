// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/massalabs/thyra/api/swagger/server/models"
)

// WebsiteCreatorUploadOKCode is the HTTP code returned for type WebsiteCreatorUploadOK
const WebsiteCreatorUploadOKCode int = 200

/*WebsiteCreatorUploadOK Website's chunk deployed.

swagger:response websiteCreatorUploadOK
*/
type WebsiteCreatorUploadOK struct {

	/*
	  In: Body
	*/
	Payload *models.Websites `json:"body,omitempty"`
}

// NewWebsiteCreatorUploadOK creates WebsiteCreatorUploadOK with default headers values
func NewWebsiteCreatorUploadOK() *WebsiteCreatorUploadOK {

	return &WebsiteCreatorUploadOK{}
}

// WithPayload adds the payload to the website creator upload o k response
func (o *WebsiteCreatorUploadOK) WithPayload(payload *models.Websites) *WebsiteCreatorUploadOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator upload o k response
func (o *WebsiteCreatorUploadOK) SetPayload(payload *models.Websites) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorUploadOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WebsiteCreatorUploadBadRequestCode is the HTTP code returned for type WebsiteCreatorUploadBadRequest
const WebsiteCreatorUploadBadRequestCode int = 400

/*WebsiteCreatorUploadBadRequest Bad request.

swagger:response websiteCreatorUploadBadRequest
*/
type WebsiteCreatorUploadBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewWebsiteCreatorUploadBadRequest creates WebsiteCreatorUploadBadRequest with default headers values
func NewWebsiteCreatorUploadBadRequest() *WebsiteCreatorUploadBadRequest {

	return &WebsiteCreatorUploadBadRequest{}
}

// WithPayload adds the payload to the website creator upload bad request response
func (o *WebsiteCreatorUploadBadRequest) WithPayload(payload *models.Error) *WebsiteCreatorUploadBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator upload bad request response
func (o *WebsiteCreatorUploadBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorUploadBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WebsiteCreatorUploadUnprocessableEntityCode is the HTTP code returned for type WebsiteCreatorUploadUnprocessableEntity
const WebsiteCreatorUploadUnprocessableEntityCode int = 422

/*WebsiteCreatorUploadUnprocessableEntity Unprocessable Entity - syntax is correct, but the server was unable to process the contained instructions.

swagger:response websiteCreatorUploadUnprocessableEntity
*/
type WebsiteCreatorUploadUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewWebsiteCreatorUploadUnprocessableEntity creates WebsiteCreatorUploadUnprocessableEntity with default headers values
func NewWebsiteCreatorUploadUnprocessableEntity() *WebsiteCreatorUploadUnprocessableEntity {

	return &WebsiteCreatorUploadUnprocessableEntity{}
}

// WithPayload adds the payload to the website creator upload unprocessable entity response
func (o *WebsiteCreatorUploadUnprocessableEntity) WithPayload(payload *models.Error) *WebsiteCreatorUploadUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator upload unprocessable entity response
func (o *WebsiteCreatorUploadUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorUploadUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WebsiteCreatorUploadInternalServerErrorCode is the HTTP code returned for type WebsiteCreatorUploadInternalServerError
const WebsiteCreatorUploadInternalServerErrorCode int = 500

/*WebsiteCreatorUploadInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response websiteCreatorUploadInternalServerError
*/
type WebsiteCreatorUploadInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewWebsiteCreatorUploadInternalServerError creates WebsiteCreatorUploadInternalServerError with default headers values
func NewWebsiteCreatorUploadInternalServerError() *WebsiteCreatorUploadInternalServerError {

	return &WebsiteCreatorUploadInternalServerError{}
}

// WithPayload adds the payload to the website creator upload internal server error response
func (o *WebsiteCreatorUploadInternalServerError) WithPayload(payload *models.Error) *WebsiteCreatorUploadInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator upload internal server error response
func (o *WebsiteCreatorUploadInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorUploadInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
