// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/massalabs/thyra/api/swagger/server/models"
)

// ThyraEventsGetterOKCode is the HTTP code returned for type ThyraEventsGetterOK
const ThyraEventsGetterOKCode int = 200

/*ThyraEventsGetterOK Event retrieved

swagger:response thyraEventsGetterOK
*/
type ThyraEventsGetterOK struct {
}

// NewThyraEventsGetterOK creates ThyraEventsGetterOK with default headers values
func NewThyraEventsGetterOK() *ThyraEventsGetterOK {

	return &ThyraEventsGetterOK{}
}

// WriteResponse to the client
func (o *ThyraEventsGetterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// ThyraEventsGetterBadRequestCode is the HTTP code returned for type ThyraEventsGetterBadRequest
const ThyraEventsGetterBadRequestCode int = 400

/*ThyraEventsGetterBadRequest Bad request.

swagger:response thyraEventsGetterBadRequest
*/
type ThyraEventsGetterBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewThyraEventsGetterBadRequest creates ThyraEventsGetterBadRequest with default headers values
func NewThyraEventsGetterBadRequest() *ThyraEventsGetterBadRequest {

	return &ThyraEventsGetterBadRequest{}
}

// WithPayload adds the payload to the thyra events getter bad request response
func (o *ThyraEventsGetterBadRequest) WithPayload(payload *models.Error) *ThyraEventsGetterBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thyra events getter bad request response
func (o *ThyraEventsGetterBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThyraEventsGetterBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ThyraEventsGetterInternalServerErrorCode is the HTTP code returned for type ThyraEventsGetterInternalServerError
const ThyraEventsGetterInternalServerErrorCode int = 500

/*ThyraEventsGetterInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response thyraEventsGetterInternalServerError
*/
type ThyraEventsGetterInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewThyraEventsGetterInternalServerError creates ThyraEventsGetterInternalServerError with default headers values
func NewThyraEventsGetterInternalServerError() *ThyraEventsGetterInternalServerError {

	return &ThyraEventsGetterInternalServerError{}
}

// WithPayload adds the payload to the thyra events getter internal server error response
func (o *ThyraEventsGetterInternalServerError) WithPayload(payload *models.Error) *ThyraEventsGetterInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the thyra events getter internal server error response
func (o *ThyraEventsGetterInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ThyraEventsGetterInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
