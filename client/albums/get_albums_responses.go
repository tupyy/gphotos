// Code generated by go-swagger; DO NOT EDIT.

package albums

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/tupyy/gophoto/models"
)

// GetAlbumsReader is a Reader for the GetAlbums structure.
type GetAlbumsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAlbumsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAlbumsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetAlbumsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetAlbumsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAlbumsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 503:
		result := NewGetAlbumsServiceUnavailable()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetAlbumsOK creates a GetAlbumsOK with default headers values
func NewGetAlbumsOK() *GetAlbumsOK {
	return &GetAlbumsOK{}
}

/* GetAlbumsOK describes a response with status code 200, with default header values.

list of albums that can be accessed by the logged user
*/
type GetAlbumsOK struct {
	Payload *models.Albums
}

// IsSuccess returns true when this get albums o k response has a 2xx status code
func (o *GetAlbumsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get albums o k response has a 3xx status code
func (o *GetAlbumsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get albums o k response has a 4xx status code
func (o *GetAlbumsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get albums o k response has a 5xx status code
func (o *GetAlbumsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get albums o k response a status code equal to that given
func (o *GetAlbumsOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetAlbumsOK) Error() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsOK  %+v", 200, o.Payload)
}

func (o *GetAlbumsOK) String() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsOK  %+v", 200, o.Payload)
}

func (o *GetAlbumsOK) GetPayload() *models.Albums {
	return o.Payload
}

func (o *GetAlbumsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Albums)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAlbumsUnauthorized creates a GetAlbumsUnauthorized with default headers values
func NewGetAlbumsUnauthorized() *GetAlbumsUnauthorized {
	return &GetAlbumsUnauthorized{}
}

/* GetAlbumsUnauthorized describes a response with status code 401, with default header values.

Not authenticated.
*/
type GetAlbumsUnauthorized struct {
}

// IsSuccess returns true when this get albums unauthorized response has a 2xx status code
func (o *GetAlbumsUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get albums unauthorized response has a 3xx status code
func (o *GetAlbumsUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get albums unauthorized response has a 4xx status code
func (o *GetAlbumsUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get albums unauthorized response has a 5xx status code
func (o *GetAlbumsUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get albums unauthorized response a status code equal to that given
func (o *GetAlbumsUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *GetAlbumsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsUnauthorized ", 401)
}

func (o *GetAlbumsUnauthorized) String() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsUnauthorized ", 401)
}

func (o *GetAlbumsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetAlbumsForbidden creates a GetAlbumsForbidden with default headers values
func NewGetAlbumsForbidden() *GetAlbumsForbidden {
	return &GetAlbumsForbidden{}
}

/* GetAlbumsForbidden describes a response with status code 403, with default header values.

Forbidden.
*/
type GetAlbumsForbidden struct {
}

// IsSuccess returns true when this get albums forbidden response has a 2xx status code
func (o *GetAlbumsForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get albums forbidden response has a 3xx status code
func (o *GetAlbumsForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get albums forbidden response has a 4xx status code
func (o *GetAlbumsForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this get albums forbidden response has a 5xx status code
func (o *GetAlbumsForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this get albums forbidden response a status code equal to that given
func (o *GetAlbumsForbidden) IsCode(code int) bool {
	return code == 403
}

func (o *GetAlbumsForbidden) Error() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsForbidden ", 403)
}

func (o *GetAlbumsForbidden) String() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsForbidden ", 403)
}

func (o *GetAlbumsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetAlbumsInternalServerError creates a GetAlbumsInternalServerError with default headers values
func NewGetAlbumsInternalServerError() *GetAlbumsInternalServerError {
	return &GetAlbumsInternalServerError{}
}

/* GetAlbumsInternalServerError describes a response with status code 500, with default header values.

Internal error.
*/
type GetAlbumsInternalServerError struct {
}

// IsSuccess returns true when this get albums internal server error response has a 2xx status code
func (o *GetAlbumsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get albums internal server error response has a 3xx status code
func (o *GetAlbumsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get albums internal server error response has a 4xx status code
func (o *GetAlbumsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get albums internal server error response has a 5xx status code
func (o *GetAlbumsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get albums internal server error response a status code equal to that given
func (o *GetAlbumsInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *GetAlbumsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsInternalServerError ", 500)
}

func (o *GetAlbumsInternalServerError) String() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsInternalServerError ", 500)
}

func (o *GetAlbumsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetAlbumsServiceUnavailable creates a GetAlbumsServiceUnavailable with default headers values
func NewGetAlbumsServiceUnavailable() *GetAlbumsServiceUnavailable {
	return &GetAlbumsServiceUnavailable{}
}

/* GetAlbumsServiceUnavailable describes a response with status code 503, with default header values.

Not available.
*/
type GetAlbumsServiceUnavailable struct {
}

// IsSuccess returns true when this get albums service unavailable response has a 2xx status code
func (o *GetAlbumsServiceUnavailable) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get albums service unavailable response has a 3xx status code
func (o *GetAlbumsServiceUnavailable) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get albums service unavailable response has a 4xx status code
func (o *GetAlbumsServiceUnavailable) IsClientError() bool {
	return false
}

// IsServerError returns true when this get albums service unavailable response has a 5xx status code
func (o *GetAlbumsServiceUnavailable) IsServerError() bool {
	return true
}

// IsCode returns true when this get albums service unavailable response a status code equal to that given
func (o *GetAlbumsServiceUnavailable) IsCode(code int) bool {
	return code == 503
}

func (o *GetAlbumsServiceUnavailable) Error() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsServiceUnavailable ", 503)
}

func (o *GetAlbumsServiceUnavailable) String() string {
	return fmt.Sprintf("[GET /api/gphotos/v1/albums][%d] getAlbumsServiceUnavailable ", 503)
}

func (o *GetAlbumsServiceUnavailable) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
