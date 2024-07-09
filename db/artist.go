package db

import (
	"context"
	"strconv"
	"time"

	"github.com/Fan-Fuse/artist-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Artist struct {
	gorm.Model
	Name        string
	SpotifyID   string
	Images      []Image `gorm:"polymorphic:Owner;"`
	Albums      []Album `gorm:"many2many:artist_albums;"` // many to many relationship
	LastUpdated time.Time
}

func CreateArtist(ctx context.Context, in *proto.Artist) (*Artist, error) {
	// Serialize the proto Artist to a db Artist
	var images []Image
	for _, img := range in.Images {
		images = append(images, Image{
			URL:    img.Url,
			Width:  img.Width,
			Height: img.Height,
		})
	}

	// TODO: Relocate this logic to album.go
	var albums []Album
	for _, alb := range in.Albums {
		var albumImages []Image
		for _, img := range alb.Images {
			albumImages = append(albumImages, Image{
				URL:    img.Url,
				Width:  img.Width,
				Height: img.Height,
			})
		}

		albums = append(albums, Album{
			Name:        alb.Name,
			SpotifyID:   alb.Externals.Spotify,
			ReleaseDate: alb.ReleaseDate.AsTime(),
			Images:      albumImages,
		})
	}

	artist := Artist{
		Name:        in.Name,
		SpotifyID:   in.Externals.Spotify,
		LastUpdated: time.Now(),
		Images:      images,
		Albums:      albums,
	}

	// TODO: Test if this creates the artist with the albums and images
	if err := DB.Create(&artist).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

func GetArtistByID(id string) (*proto.Artist, error) {
	var artist Artist
	if err := DB.Preload("Albums").Preload("Images").First(&artist, id).Error; err != nil {
		return nil, err
	}
	// Serialize the db Artist to a proto Artist
	var images []*proto.Image
	for _, img := range artist.Images {
		images = append(images, &proto.Image{
			Url:    img.URL,
			Width:  img.Width,
			Height: img.Height,
		})
	}

	var albums []*proto.Album
	for _, alb := range artist.Albums {
		var albumImages []*proto.Image
		for _, img := range alb.Images {
			albumImages = append(albumImages, &proto.Image{
				Url:    img.URL,
				Width:  img.Width,
				Height: img.Height,
			})
		}

		albums = append(albums, &proto.Album{
			Name:        alb.Name,
			ReleaseDate: &timestamppb.Timestamp{Seconds: alb.ReleaseDate.Unix()},
			Images:      albumImages,
		})
	}

	return &proto.Artist{
		Id:     strconv.Itoa(int(artist.ID)),
		Name:   artist.Name,
		Images: images,
		Albums: albums,
		Externals: &proto.Externals{
			Spotify: artist.SpotifyID,
		},
	}, nil
}

func GetArtistBySpotifyID(spotifyID string) (*Artist, error) {
	var artist Artist
	if err := DB.Where("spotify_id = ?", spotifyID).First(&artist).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

func GetArtists(limit, offset int32) ([]*proto.Artist, error) {
	var artists []Artist
	if err := DB.Offset(int(offset)).Limit(int(limit)).Find(&artists).Error; err != nil {
		return nil, err
	}

	var protoArtists []*proto.Artist
	for _, artist := range artists {
		var images []*proto.Image
		for _, img := range artist.Images {
			images = append(images, &proto.Image{
				Url:    img.URL,
				Width:  img.Width,
				Height: img.Height,
			})
		}

		// We won't load the albums here, as it would be too much data
		protoArtists = append(protoArtists, &proto.Artist{
			Id:     strconv.Itoa(int(artist.ID)),
			Name:   artist.Name,
			Images: images,
			Externals: &proto.Externals{
				Spotify: artist.SpotifyID,
			},
		})
	}

	return protoArtists, nil
}

func AddAlbumToArtist(artist *Artist, album *Album) error {
	artist.Albums = append(artist.Albums, *album)
	if err := DB.Save(&artist).Error; err != nil {
		return err
	}
	return nil
}

func FilterArtists(filter *proto.ArtistFilter) ([]*proto.Artist, error) {
	var artists []Artist
	// Decide which fields to filter on
	if filter.Name != "" {
		DB = DB.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.LastUpdated != nil {
		if filter.LastUpdatedGt {
			DB = DB.Where("last_updated > ?", filter.LastUpdated.AsTime())
		} else {
			DB = DB.Where("last_updated < ?", filter.LastUpdated.AsTime())
		}
	}

	// Retrieve the artists
	if err := DB.Find(&artists).Error; err != nil {
		return nil, err
	}

	var protoArtists []*proto.Artist
	for _, artist := range artists {
		var images []*proto.Image
		for _, img := range artist.Images {
			images = append(images, &proto.Image{
				Url:    img.URL,
				Width:  img.Width,
				Height: img.Height,
			})
		}

		// We won't load the albums here, as it would be too much data
		protoArtists = append(protoArtists, &proto.Artist{
			Id:     strconv.Itoa(int(artist.ID)),
			Name:   artist.Name,
			Images: images,
			Externals: &proto.Externals{
				Spotify: artist.SpotifyID,
			},
		})
	}

	return protoArtists, nil
}
