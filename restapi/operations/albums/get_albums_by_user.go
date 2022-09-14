// Code generated by go-swagger; DO NOT EDIT.

package albums

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetAlbumsByUserHandlerFunc turns a function with the right signature into a get albums by user handler
type GetAlbumsByUserHandlerFunc func(GetAlbumsByUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAlbumsByUserHandlerFunc) Handle(params GetAlbumsByUserParams) middleware.Responder {
	return fn(params)
}

// GetAlbumsByUserHandler interface for that can handle valid get albums by user params
type GetAlbumsByUserHandler interface {
	Handle(GetAlbumsByUserParams) middleware.Responder
}

// NewGetAlbumsByUser creates a new http.Handler for the get albums by user operation
func NewGetAlbumsByUser(ctx *middleware.Context, handler GetAlbumsByUserHandler) *GetAlbumsByUser {
	return &GetAlbumsByUser{Context: ctx, Handler: handler}
}

/* GetAlbumsByUser swagger:route GET /api/gphotos/v1/albums/users/{user_id} Albums getAlbumsByUser

get all user's album with the logged user can access

*/
type GetAlbumsByUser struct {
	Context *middleware.Context
	Handler GetAlbumsByUserHandler
}

func (o *GetAlbumsByUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAlbumsByUserParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
