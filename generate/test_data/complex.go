package test_data

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

import "time"

type Complex struct {
	Simple

	Caption string `json:"caption"`
	Credit  string `json:"credit"`
	Crops   []struct {
		Height       float64 `json:"height"`
		Name         string  `json:"name"`
		Path         string  `json:"path" description:"full path to the cropped image file"`
		RelativePath string  `json:"relativePath" description:"a long"`
		Width        float64 `json:"width"`
	} `json:"crops"`
	Cutline        string    `json:"cutline,omitempty"`
	DatePhotoTaken time.Time `json:"datePhotoTaken"`
	Orientation    string    `json:"orientation"`
	OriginalSize   struct {
		Height float64 `json:"height"`
		Width  float64 `json:"width"`
	} `json:"originalSize"`
	Type string `json:"type" description:"a type"`
	URL  struct {
		Absolute string `json:"absolute"`
		Publish  string `json:"publish"`
	} `json:"URL"`
}
