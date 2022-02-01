package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"

	"github.com/taehoio/oneonone/config"
	"github.com/taehoio/oneonone/server"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger := logrus.StandardLogger()

	setting := config.NewSetting()
	cfg := config.NewConfig(setting, logger)

	if err := runServer(cfg); err != nil {
		logger.Fatal(err)
	}
}

func runServer(cfg config.Config) error {
	log := cfg.Logger()

	if cfg.Setting().ShouldProfile {
		if err := setUpProfiler(cfg.Setting().ServiceName); err != nil {
			return err
		}
	}

	if cfg.Setting().ShouldTrace {
		tp, err := setUpTracing(cfg.Setting().ServiceName)
		if err != nil {
			return err
		}
		defer tp.ForceFlush(context.Background())
	}

	grpcServer, err := server.NewGRPCServer(cfg)
	if err != nil {
		return err
	}

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.Setting().GRPCServerPort)
		if err != nil {
			log.Fatal(err)
		}

		log.WithField("port", cfg.Setting().GRPCServerPort).Infof("Starting %s gRPC server", cfg.Setting().ServiceName)
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit

	time.Sleep(time.Duration(cfg.Setting().GracefulShutdownTimeoutMs) * time.Millisecond)

	log.Infof("Stopping %s gRPC server", cfg.Setting().ServiceName)
	grpcServer.GracefulStop()

	return nil
}

func setUpProfiler(serviceName string) error {
	pc := profiler.Config{
		Service: serviceName,
	}
	if err := profiler.Start(pc); err != nil {
		return err
	}
	return nil
}

func setUpTracing(serviceName string) (*trace.TracerProvider, error) {
	exporter, err := texporter.New()
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}
