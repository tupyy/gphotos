// Code generated by go-swagger; DO NOT EDIT.

package albums

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetAlbumsByIDHandlerFunc turns a function with the right signature into a get albums by ID handler
type GetAlbumsByIDHandlerFunc func(GetAlbumsByIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAlbumsByIDHandlerFunc) Handle(params GetAlbumsByIDParams) middleware.Responder {
	return fn(params)
}

// GetAlbumsByIDHandler interface for that can handle valid get albums by ID params
type GetAlbumsByIDHandler interface {
	Handle(GetAlbumsByIDParams) middleware.Responder
}

// NewGetAlbumsByID creates a new http.Handler for the get albums by ID operation
func NewGetAlbumsByID(ctx *middleware.Context, handler GetAlbumsByIDHandler) *GetAlbumsByID {
	return &GetAlbumsByID{Context: ctx, Handler: handler}
}

/* GetAlbumsByID swagger:route GET /api/gphotos/v1/albums/{album_id} Albums getAlbumsById

get album by id

*/
type GetAlbumsByID struct {
	Context *middleware.Context
	Handler GetAlbumsByIDHandler
}

func (o *GetAlbumsByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAlbumsByIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
