package nonest

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

import "time"

type Simple struct {
	Height      int64              `json:"height,omitempty"`
	SomeDateObj *SimpleSomeDateObj `json:"someDateObj,omitempty"`
	Type        string             `json:"type"`
	Visible     bool               `json:"visible,omitempty"`
	Width       float64            `json:"width,omitempty"`
}

type SimpleSomeDateObj struct {
	Dates []time.Time `json:"dates,omitempty"`
}
