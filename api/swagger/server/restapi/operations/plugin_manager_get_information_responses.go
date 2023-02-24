// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/massalabs/thyra/api/swagger/server/models"
)

// PluginManagerGetInformationOKCode is the HTTP code returned for type PluginManagerGetInformationOK
const PluginManagerGetInformationOKCode int = 200

/*
PluginManagerGetInformationOK Get execution information from the plugin.

swagger:response pluginManagerGetInformationOK
*/
type PluginManagerGetInformationOK struct {

	/*
	  In: Body
	*/
	Payload *PluginManagerGetInformationOKBody `json:"body,omitempty"`
}

// NewPluginManagerGetInformationOK creates PluginManagerGetInformationOK with default headers values
func NewPluginManagerGetInformationOK() *PluginManagerGetInformationOK {

	return &PluginManagerGetInformationOK{}
}

// WithPayload adds the payload to the plugin manager get information o k response
func (o *PluginManagerGetInformationOK) WithPayload(payload *PluginManagerGetInformationOKBody) *PluginManagerGetInformationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the plugin manager get information o k response
func (o *PluginManagerGetInformationOK) SetPayload(payload *PluginManagerGetInformationOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PluginManagerGetInformationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PluginManagerGetInformationBadRequestCode is the HTTP code returned for type PluginManagerGetInformationBadRequest
const PluginManagerGetInformationBadRequestCode int = 400

/*
PluginManagerGetInformationBadRequest Bad request.

swagger:response pluginManagerGetInformationBadRequest
*/
type PluginManagerGetInformationBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPluginManagerGetInformationBadRequest creates PluginManagerGetInformationBadRequest with default headers values
func NewPluginManagerGetInformationBadRequest() *PluginManagerGetInformationBadRequest {

	return &PluginManagerGetInformationBadRequest{}
}

// WithPayload adds the payload to the plugin manager get information bad request response
func (o *PluginManagerGetInformationBadRequest) WithPayload(payload *models.Error) *PluginManagerGetInformationBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the plugin manager get information bad request response
func (o *PluginManagerGetInformationBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PluginManagerGetInformationBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PluginManagerGetInformationNotFoundCode is the HTTP code returned for type PluginManagerGetInformationNotFound
const PluginManagerGetInformationNotFoundCode int = 404

/*
PluginManagerGetInformationNotFound Not found.

swagger:response pluginManagerGetInformationNotFound
*/
type PluginManagerGetInformationNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPluginManagerGetInformationNotFound creates PluginManagerGetInformationNotFound with default headers values
func NewPluginManagerGetInformationNotFound() *PluginManagerGetInformationNotFound {

	return &PluginManagerGetInformationNotFound{}
}

// WithPayload adds the payload to the plugin manager get information not found response
func (o *PluginManagerGetInformationNotFound) WithPayload(payload *models.Error) *PluginManagerGetInformationNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the plugin manager get information not found response
func (o *PluginManagerGetInformationNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PluginManagerGetInformationNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PluginManagerGetInformationUnprocessableEntityCode is the HTTP code returned for type PluginManagerGetInformationUnprocessableEntity
const PluginManagerGetInformationUnprocessableEntityCode int = 422

/*
PluginManagerGetInformationUnprocessableEntity Unprocessable Entity - The syntax is correct, but the server was unable to process the contained instructions.

swagger:response pluginManagerGetInformationUnprocessableEntity
*/
type PluginManagerGetInformationUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPluginManagerGetInformationUnprocessableEntity creates PluginManagerGetInformationUnprocessableEntity with default headers values
func NewPluginManagerGetInformationUnprocessableEntity() *PluginManagerGetInformationUnprocessableEntity {

	return &PluginManagerGetInformationUnprocessableEntity{}
}

// WithPayload adds the payload to the plugin manager get information unprocessable entity response
func (o *PluginManagerGetInformationUnprocessableEntity) WithPayload(payload *models.Error) *PluginManagerGetInformationUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the plugin manager get information unprocessable entity response
func (o *PluginManagerGetInformationUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PluginManagerGetInformationUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PluginManagerGetInformationInternalServerErrorCode is the HTTP code returned for type PluginManagerGetInformationInternalServerError
const PluginManagerGetInformationInternalServerErrorCode int = 500

/*
PluginManagerGetInformationInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response pluginManagerGetInformationInternalServerError
*/
type PluginManagerGetInformationInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPluginManagerGetInformationInternalServerError creates PluginManagerGetInformationInternalServerError with default headers values
func NewPluginManagerGetInformationInternalServerError() *PluginManagerGetInformationInternalServerError {

	return &PluginManagerGetInformationInternalServerError{}
}

// WithPayload adds the payload to the plugin manager get information internal server error response
func (o *PluginManagerGetInformationInternalServerError) WithPayload(payload *models.Error) *PluginManagerGetInformationInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the plugin manager get information internal server error response
func (o *PluginManagerGetInformationInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PluginManagerGetInformationInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
