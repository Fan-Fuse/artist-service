package service

import (
	"github.com/Fan-Fuse/artist-service/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedArtistServiceServer
}

// RegisterServer registers the server with the gRPC server
func RegisterServer(s *grpc.Server) {
	proto.RegisterArtistServiceServer(s, &server{})
}
