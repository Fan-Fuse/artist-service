package db

import (
	"time"

	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Name        string
	SpotifyID   string
	ReleaseDate time.Time
	ImageURL    string
	Artists     []Artist `gorm:"many2many:artist_albums;"` // many to many relationship
}

func CreateAlbum(name, spotifyID, imageUrl string, releaseDate time.Time) (*Album, error) {
	album := Album{
		Name:        name,
		SpotifyID:   spotifyID,
		ImageURL:    imageUrl,
		ReleaseDate: releaseDate,
	}
	if err := DB.Create(&album).Error; err != nil {
		return nil, err
	}
	return &album, nil
}

func GetAlbumByID(id uint) (*Album, error) {
	var album Album
	if err := DB.Preload("Artists").First(&album, id).Error; err != nil {
		return nil, err
	}
	return &album, nil
}

func GetAlbumBySpotifyID(spotifyID string) (*Album, error) {
	var album Album
	if err := DB.Where("spotify_id = ?", spotifyID).First(&album).Error; err != nil {
		return nil, err
	}
	return &album, nil
}
