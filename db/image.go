package db

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	URL    string
	Width  int32
	Height int32
}
