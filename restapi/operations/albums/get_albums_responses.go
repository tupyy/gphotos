// Code generated by go-swagger; DO NOT EDIT.

package albums

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/tupyy/gophoto/models"
)

// GetAlbumsOKCode is the HTTP code returned for type GetAlbumsOK
const GetAlbumsOKCode int = 200

/*GetAlbumsOK list of albums that can be accessed by the logged user

swagger:response getAlbumsOK
*/
type GetAlbumsOK struct {

	/*
	  In: Body
	*/
	Payload *models.Albums `json:"body,omitempty"`
}

// NewGetAlbumsOK creates GetAlbumsOK with default headers values
func NewGetAlbumsOK() *GetAlbumsOK {

	return &GetAlbumsOK{}
}

// WithPayload adds the payload to the get albums o k response
func (o *GetAlbumsOK) WithPayload(payload *models.Albums) *GetAlbumsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get albums o k response
func (o *GetAlbumsOK) SetPayload(payload *models.Albums) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAlbumsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAlbumsUnauthorizedCode is the HTTP code returned for type GetAlbumsUnauthorized
const GetAlbumsUnauthorizedCode int = 401

/*GetAlbumsUnauthorized Not authenticated.

swagger:response getAlbumsUnauthorized
*/
type GetAlbumsUnauthorized struct {
}

// NewGetAlbumsUnauthorized creates GetAlbumsUnauthorized with default headers values
func NewGetAlbumsUnauthorized() *GetAlbumsUnauthorized {

	return &GetAlbumsUnauthorized{}
}

// WriteResponse to the client
func (o *GetAlbumsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// GetAlbumsForbiddenCode is the HTTP code returned for type GetAlbumsForbidden
const GetAlbumsForbiddenCode int = 403

/*GetAlbumsForbidden Forbidden.

swagger:response getAlbumsForbidden
*/
type GetAlbumsForbidden struct {
}

// NewGetAlbumsForbidden creates GetAlbumsForbidden with default headers values
func NewGetAlbumsForbidden() *GetAlbumsForbidden {

	return &GetAlbumsForbidden{}
}

// WriteResponse to the client
func (o *GetAlbumsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// GetAlbumsInternalServerErrorCode is the HTTP code returned for type GetAlbumsInternalServerError
const GetAlbumsInternalServerErrorCode int = 500

/*GetAlbumsInternalServerError Internal error.

swagger:response getAlbumsInternalServerError
*/
type GetAlbumsInternalServerError struct {
}

// NewGetAlbumsInternalServerError creates GetAlbumsInternalServerError with default headers values
func NewGetAlbumsInternalServerError() *GetAlbumsInternalServerError {

	return &GetAlbumsInternalServerError{}
}

// WriteResponse to the client
func (o *GetAlbumsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}

// GetAlbumsServiceUnavailableCode is the HTTP code returned for type GetAlbumsServiceUnavailable
const GetAlbumsServiceUnavailableCode int = 503

/*GetAlbumsServiceUnavailable Not available.

swagger:response getAlbumsServiceUnavailable
*/
type GetAlbumsServiceUnavailable struct {
}

// NewGetAlbumsServiceUnavailable creates GetAlbumsServiceUnavailable with default headers values
func NewGetAlbumsServiceUnavailable() *GetAlbumsServiceUnavailable {

	return &GetAlbumsServiceUnavailable{}
}

// WriteResponse to the client
func (o *GetAlbumsServiceUnavailable) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(503)
}
