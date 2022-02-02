package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
	"github.com/taehoio/oneonone/config"
	"github.com/taehoio/oneonone/server/handler"
)

type OneononeServiceServer struct {
	oneononev1.OneononeServiceServer

	cfg config.Config
	db  *sql.DB
}

func NewOneononeServiceServer(cfg config.Config) (*OneononeServiceServer, error) {
	db, err := newMySQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &OneononeServiceServer{
		cfg: cfg,
		db:  db,
	}, nil
}

func newMySQLDB(cfg config.Config) (*sql.DB, error) {
	mysqlCfg := mysql.Config{
		Net:                  cfg.Setting().MysqlNetworkType,
		Addr:                 cfg.Setting().MysqlAddress,
		User:                 cfg.Setting().MysqlUser,
		Passwd:               cfg.Setting().MysqlPassword,
		DBName:               cfg.Setting().MysqlDatabaseName,
		AllowNativePasswords: true,
		ParseTime:            true,
		TLSConfig:            "preferred",
	}

	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *OneononeServiceServer) HealthCheck(ctx context.Context, req *oneononev1.HealthCheckRequest) (*oneononev1.HealthCheckResponse, error) {
	return &oneononev1.HealthCheckResponse{}, nil
}

func (s *OneononeServiceServer) ListCategories(ctx context.Context, req *oneononev1.ListCategoriesRequest) (*oneononev1.ListCategoriesResponse, error) {
	return handler.ListCategories(s.db, &oneononeddlv1.Category{})(ctx, req)
}

func (s *OneononeServiceServer) ListQuestionsByCategoryId(ctx context.Context, req *oneononev1.ListQuestionsByCategoryIdRequest) (*oneononev1.ListQuestionsByCategoryIdResponse, error) {
	return handler.ListQuestionsByCategoryId(s.db, &oneononeddlv1.CategoryQuestion{}, &oneononeddlv1.Question{})(ctx, req)
}

func (s *OneononeServiceServer) ListQuestiGetRandomQuestiononsByCategoryId(ctx context.Context, req *oneononev1.GetRandomQuestionRequest) (*oneononev1.GetRandomQuestionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *OneononeServiceServer) GetRandomQuestionByCategoryId(ctx context.Context, req *oneononev1.GetRandomQuestionByCategoryIdRequest) (*oneononev1.GetRandomQuestionByCategoryIdResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func NewGRPCServer(cfg config.Config) (*grpc.Server, error) {
	logrus.ErrorKey = "grpc.error"
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			otelgrpc.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(
					grpc_ctxtags.CodeGenRequestFieldExtractor,
				),
			),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 30 * time.Second,
		}),
	)

	oneononeServiceServer, err := NewOneononeServiceServer(cfg)
	if err != nil {
		return nil, err
	}

	oneononev1.RegisterOneononeServiceServer(grpcServer, oneononeServiceServer)
	reflection.Register(grpcServer)

	return grpcServer, nil
}
