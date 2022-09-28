// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/massalabs/thyra/api/swagger/server/models"
)

// WebsiteCreatorPrepareOKCode is the HTTP code returned for type WebsiteCreatorPrepareOK
const WebsiteCreatorPrepareOKCode int = 200

/*WebsiteCreatorPrepareOK New website created.

swagger:response websiteCreatorPrepareOK
*/
type WebsiteCreatorPrepareOK struct {

	/*
	  In: Body
	*/
	Payload *models.Websites `json:"body,omitempty"`
}

// NewWebsiteCreatorPrepareOK creates WebsiteCreatorPrepareOK with default headers values
func NewWebsiteCreatorPrepareOK() *WebsiteCreatorPrepareOK {

	return &WebsiteCreatorPrepareOK{}
}

// WithPayload adds the payload to the website creator prepare o k response
func (o *WebsiteCreatorPrepareOK) WithPayload(payload *models.Websites) *WebsiteCreatorPrepareOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator prepare o k response
func (o *WebsiteCreatorPrepareOK) SetPayload(payload *models.Websites) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorPrepareOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WebsiteCreatorPrepareBadRequestCode is the HTTP code returned for type WebsiteCreatorPrepareBadRequest
const WebsiteCreatorPrepareBadRequestCode int = 400

/*WebsiteCreatorPrepareBadRequest Bad request.

swagger:response websiteCreatorPrepareBadRequest
*/
type WebsiteCreatorPrepareBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewWebsiteCreatorPrepareBadRequest creates WebsiteCreatorPrepareBadRequest with default headers values
func NewWebsiteCreatorPrepareBadRequest() *WebsiteCreatorPrepareBadRequest {

	return &WebsiteCreatorPrepareBadRequest{}
}

// WithPayload adds the payload to the website creator prepare bad request response
func (o *WebsiteCreatorPrepareBadRequest) WithPayload(payload *models.Error) *WebsiteCreatorPrepareBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator prepare bad request response
func (o *WebsiteCreatorPrepareBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorPrepareBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WebsiteCreatorPrepareUnprocessableEntityCode is the HTTP code returned for type WebsiteCreatorPrepareUnprocessableEntity
const WebsiteCreatorPrepareUnprocessableEntityCode int = 422

/*WebsiteCreatorPrepareUnprocessableEntity Unprocessable Entity - syntax is correct, but the server was unable to process the contained instructions.

swagger:response websiteCreatorPrepareUnprocessableEntity
*/
type WebsiteCreatorPrepareUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewWebsiteCreatorPrepareUnprocessableEntity creates WebsiteCreatorPrepareUnprocessableEntity with default headers values
func NewWebsiteCreatorPrepareUnprocessableEntity() *WebsiteCreatorPrepareUnprocessableEntity {

	return &WebsiteCreatorPrepareUnprocessableEntity{}
}

// WithPayload adds the payload to the website creator prepare unprocessable entity response
func (o *WebsiteCreatorPrepareUnprocessableEntity) WithPayload(payload *models.Error) *WebsiteCreatorPrepareUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator prepare unprocessable entity response
func (o *WebsiteCreatorPrepareUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorPrepareUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WebsiteCreatorPrepareInternalServerErrorCode is the HTTP code returned for type WebsiteCreatorPrepareInternalServerError
const WebsiteCreatorPrepareInternalServerErrorCode int = 500

/*WebsiteCreatorPrepareInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response websiteCreatorPrepareInternalServerError
*/
type WebsiteCreatorPrepareInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewWebsiteCreatorPrepareInternalServerError creates WebsiteCreatorPrepareInternalServerError with default headers values
func NewWebsiteCreatorPrepareInternalServerError() *WebsiteCreatorPrepareInternalServerError {

	return &WebsiteCreatorPrepareInternalServerError{}
}

// WithPayload adds the payload to the website creator prepare internal server error response
func (o *WebsiteCreatorPrepareInternalServerError) WithPayload(payload *models.Error) *WebsiteCreatorPrepareInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the website creator prepare internal server error response
func (o *WebsiteCreatorPrepareInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WebsiteCreatorPrepareInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
