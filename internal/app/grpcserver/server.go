package grpcserver

import (
	"context"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/grpcserver/proto"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/service"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"github.com/carinfinin/shortener-url/internal/app/storage/storefile"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

// ShortenerServer grpc server.
type ShortenerServer struct {
	proto.UnimplementedShortenerServer
	Store   service.Repository
	Service *service.Service
	config  *config.Config
}

// New конструктор для ShortenerServer принимает кофиг.
func New(config *config.Config) (*ShortenerServer, error) {

	var server ShortenerServer

	switch {
	case config.DBPath != "":
		s, err := storepg.New(config)
		if err != nil {
			return nil, err
		}
		s.CreateTableForDB(context.Background())
		server.Store = s

	case config.FilePath != "":
		s, err := storefile.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s

	default:
		s, err := store.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s
	}

	s := service.New(server.Store, config)

	server.Service = s
	server.config = config

	return &server, nil
}

// CreateURL зосдание сокращённого url.
func (s *ShortenerServer) CreateURL(ctx context.Context, req *proto.UrlRequest) (*proto.UrlResponse, error) {
	var resp proto.UrlResponse
	short, err := s.Service.CreateURL(ctx, req.Url)
	if err != nil {
		return nil, err
	}
	resp.Result = short
	return &resp, nil
}

// CreateURLHandle создание сокращённого url.
func (s *ShortenerServer) CreateURLHandle(ctx context.Context, req *proto.UrlRequest) (*proto.UrlResponse, error) {
	var resp proto.UrlResponse
	short, err := s.Service.CreateURL(ctx, req.Url)
	if err != nil {
		return nil, err
	}
	resp.Result = short
	return &resp, nil
}

// CreateURLBatch создание пачки сокращённых урлов.
func (s *ShortenerServer) CreateURLBatch(ctx context.Context, req *proto.UrlRequestBatch) (*proto.UrlResponseBatch, error) {
	var resp proto.UrlResponseBatch

	data := make([]models.RequestBatch, 0)
	for _, v := range req.Urls {
		data = append(data, models.RequestBatch{
			ID:      v.Id,
			LongURL: v.Url,
		})
	}

	result, err := s.Service.JSONHandleBatch(ctx, data)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		resp.Urls = append(resp.Urls, &proto.Url{
			Id:  v.ID,
			Url: v.ShortURL,
		})
	}
	return &resp, nil
}

// PingDB проверка бд.
func (s *ShortenerServer) PingDB(ctx context.Context, empty *emptypb.Empty) (*proto.DBResponse, error) {

	var ping proto.DBResponse
	err := s.Service.PingDB(ctx)
	if err != nil {
		ping.Error = err.Error()
	}

	return &ping, nil
}

// GetURL получение урл.
func (s *ShortenerServer) GetURL(ctx context.Context, req *proto.UrlRequest) (*proto.UrlResponse, error) {
	var resp proto.UrlResponse
	url, err := s.Service.GetURL(ctx, req.Url)
	if err != nil {
		if errors.Is(err, storage.ErrDeleteURL) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	resp.Result = url
	return &resp, nil
}

// GetUserURLs получить  список урлов пользователя.
func (s *ShortenerServer) GetUserURLs(ctx context.Context, empty *emptypb.Empty) (*proto.UrlResponseBatch, error) {

	var resp proto.UrlResponseBatch

	data, err := s.Service.GetUserURLs(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		resp.Urls = append(resp.Urls, &proto.Url{
			Id:  v.ShortURL,
			Url: v.OriginalURL,
		})
	}
	return &resp, nil
}

// GetStats получить статистику.
func (s *ShortenerServer) GetStats(ctx context.Context, empty *emptypb.Empty) (*proto.StatResponse, error) {
	var resp proto.StatResponse

	data, err := s.Service.Stat(ctx)
	if err != nil {
		return nil, err
	}
	resp.User = int32(data.Users)
	resp.Url = int32(data.URLs)
	return &resp, nil
}

// DeleteUserURLs удаление урл.
func (s *ShortenerServer) DeleteUserURLs(ctx context.Context, req *proto.UrlResponseBatch) (*emptypb.Empty, error) {
	data := make([]string, 0)
	for _, v := range req.Urls {
		data = append(data, v.Url)
	}

	err := s.Service.DeleteUserURLs(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Start server grpc.
func Start(cfg *config.Config) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.Log.Fatal(err)
	}
	s := grpc.NewServer()

	sShort, err := New(cfg)
	if err != nil {
		logger.Log.Fatal(err)
	}
	proto.RegisterShortenerServer(s, sShort)

	logger.Log.Info("Сервер grpc начал работу")

	if err = s.Serve(listen); err != nil {
		logger.Log.Fatal(err)
	}
}
