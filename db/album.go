package db

import (
	"time"

	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Name        string
	AlbumType   string
	SpotifyID   string
	ReleaseDate time.Time
	Images      []Image  `gorm:"polymorphic:Owner;"`
	Artists     []Artist `gorm:"many2many:artist_albums;"` // many to many relationship
}
