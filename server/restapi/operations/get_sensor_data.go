// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetSensorDataHandlerFunc turns a function with the right signature into a get sensor data handler
type GetSensorDataHandlerFunc func(GetSensorDataParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetSensorDataHandlerFunc) Handle(params GetSensorDataParams) middleware.Responder {
	return fn(params)
}

// GetSensorDataHandler interface for that can handle valid get sensor data params
type GetSensorDataHandler interface {
	Handle(GetSensorDataParams) middleware.Responder
}

// NewGetSensorData creates a new http.Handler for the get sensor data operation
func NewGetSensorData(ctx *middleware.Context, handler GetSensorDataHandler) *GetSensorData {
	return &GetSensorData{Context: ctx, Handler: handler}
}

/*GetSensorData swagger:route GET /sensors/{sensorId} getSensorData

Returns the data of a particular sensor

*/
type GetSensorData struct {
	Context *middleware.Context
	Handler GetSensorDataHandler
}

func (o *GetSensorData) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetSensorDataParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}