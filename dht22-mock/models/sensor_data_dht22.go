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

// SensorDataDht22 sensor data dht22
// swagger:model sensorDataDht22
type SensorDataDht22 struct {

	// heat index
	// Required: true
	HeatIndex *float64 `json:"heatIndex"`

	// humidity
	// Required: true
	Humidity *float64 `json:"humidity"`

	// temperature
	// Required: true
	Temperature *float64 `json:"temperature"`
}

// Validate validates this sensor data dht22
func (m *SensorDataDht22) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHeatIndex(formats); err != nil {
		// prop
		res = append(res, err)
	}

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

func (m *SensorDataDht22) validateHeatIndex(formats strfmt.Registry) error {

	if err := validate.Required("heatIndex", "body", m.HeatIndex); err != nil {
		return err
	}

	return nil
}

func (m *SensorDataDht22) validateHumidity(formats strfmt.Registry) error {

	if err := validate.Required("humidity", "body", m.Humidity); err != nil {
		return err
	}

	return nil
}

func (m *SensorDataDht22) validateTemperature(formats strfmt.Registry) error {

	if err := validate.Required("temperature", "body", m.Temperature); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SensorDataDht22) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SensorDataDht22) UnmarshalBinary(b []byte) error {
	var res SensorDataDht22
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
