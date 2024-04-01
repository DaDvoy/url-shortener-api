package api

import (
	"context"
	"errors"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	"github.com/DaDvoy/url-shortener-api.git/internal/transport/protos/gen/go/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	SaveURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shorURL string) (string, error)
}

type serverAPI struct {
	apiv1.UnimplementedUrlShortenerServer
	apiI API
}

func Register(gRPC *grpc.Server, api API) {
	apiv1.RegisterUrlShortenerServer(gRPC, &serverAPI{apiI: api})
}

func (s *serverAPI) GetURL(
	ctx context.Context,
	req *apiv1.GetURLRequest,
) (*apiv1.GetURLResponse, error) {
	if req.GetShortURL() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty short url")
	}
	url, err := s.apiI.GetURL(ctx, req.GetShortURL())
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return nil, status.Error(codes.NotFound, "url not found")
		}
		return nil, status.Error(codes.Internal, "failed to get url")
	}
	return &apiv1.GetURLResponse{URL: url}, nil
}

func (s *serverAPI) PostURL(
	ctx context.Context,
	req *apiv1.PostURLRequest,
) (*apiv1.PostURLResponse, error) {
	if req.GetURL() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty url")
	}
	shortUrl, err := s.apiI.SaveURL(ctx, req.GetURL())
	if err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			return &apiv1.PostURLResponse{ShortURL: shortUrl}, nil
		}
		return nil, status.Error(codes.Internal, "failed to get short url")
	}
	return &apiv1.PostURLResponse{ShortURL: shortUrl}, nil
}
