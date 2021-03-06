// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SensorResponseRaw sensor response raw
// swagger:model SensorResponseRaw
type SensorResponseRaw struct {

	// humidity
	// Required: true
	Humidity *float64 `json:"humidity"`

	// temperature
	// Required: true
	Temperature *float64 `json:"temperature"`
}

// Validate validates this sensor response raw
func (m *SensorResponseRaw) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHumidity(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTemperature(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SensorResponseRaw) validateHumidity(formats strfmt.Registry) error {

	if err := validate.Required("humidity", "body", m.Humidity); err != nil {
		return err
	}

	return nil
}

func (m *SensorResponseRaw) validateTemperature(formats strfmt.Registry) error {

	if err := validate.Required("temperature", "body", m.Temperature); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SensorResponseRaw) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SensorResponseRaw) UnmarshalBinary(b []byte) error {
	var res SensorResponseRaw
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
