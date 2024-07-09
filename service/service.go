package service

import (
	"context"
	"strconv"

	"github.com/Fan-Fuse/artist-service/db"
	"github.com/Fan-Fuse/artist-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	proto.UnimplementedArtistServiceServer
}

// RegisterServer registers the server with the gRPC server
func RegisterServer(s *grpc.Server) {
	proto.RegisterArtistServiceServer(s, &server{})
}

func (s *server) CreateArtist(ctx context.Context, req *proto.Artist) (*proto.Id, error) {
	// Check if Artist already exists, if it does return the id and an error
	artist, err := db.GetArtistBySpotifyID(req.Externals.Spotify)
	if err == nil {
		return &proto.Id{Id: strconv.Itoa(int(artist.ID))}, status.Errorf(codes.AlreadyExists, "artist with spotify ID %s already exists", req.Externals.Spotify)
	}
	// Save the artist to the database
	artist, err = db.CreateArtist(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create artist: %v", err)
	}
	// Return the generated ID
	return &proto.Id{Id: strconv.Itoa(int(artist.ID))}, nil
}

func (s *server) GetArtist(ctx context.Context, req *proto.Id) (*proto.Artist, error) {
	artist, err := db.GetArtistByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "artist with ID %s not found", req.Id)
	}

	return artist, nil
}

func (s *server) GetArtists(ctx context.Context, req *emptypb.Empty) (*proto.Artists, error) {
	artists, err := db.GetArtists(100, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get artists: %v", err)
	}

	return &proto.Artists{Artists: artists}, nil
}

func (s *server) UpdateArtist(ctx context.Context, req *proto.Artist) (*emptypb.Empty, error) {
	// TODO: Retrieve the artist from the database based on the ID
	// TODO: Update the artist with the new data
	// TODO: Save the updated artist to the database
	// TODO: Return the updated artist

	return nil, nil
}

// GetReleases

func (s *server) FilterArtists(ctx context.Context, req *proto.ArtistFilter) (*proto.Artists, error) {
	artists, err := db.FilterArtists(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to filter artists: %v", err)
	}

	return &proto.Artists{Artists: artists}, nil
}
