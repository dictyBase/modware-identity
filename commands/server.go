package commands

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dictyBase/apihelpers/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/modware-identity/message/nats"
	"github.com/dictyBase/modware-identity/server"
	"github.com/dictyBase/modware-identity/storage"
	"github.com/dictyBase/modware-identity/storage/arangodb"
	"github.com/go-chi/cors"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gnats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer(c *cli.Context) error {
	var ds storage.DataSource
	if c.IsSet("is-secure") {
		a, err := arangodb.NewTLSDataSource(
			c.String("arangodb-user"),
			c.String("arangodb-pass"),
			c.String("arangodb-database"),
			c.String("arangodb-host"),
			c.String("arangodb-port"),
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("cannot connect to secured arangodb datasource %s", err.Error()),
				2,
			)
		}
		ds = a
	} else {
		a, err := arangodb.NewDataSource(
			c.String("arangodb-user"),
			c.String("arangodb-pass"),
			c.String("arangodb-database"),
			c.String("arangodb-host"),
			c.String("arangodb-port"),
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("cannot connect to arangodb datasource %s", err.Error()),
				2,
			)
		}
		ds = a
	}
	ms, err := nats.NewRequest(
		c.String("nats-host"),
		c.String("nats-port"),
		gnats.MaxReconnects(-1),
		gnats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to messaging server %s", err.Error()),
			2,
		)
	}
	grpcS := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(getLogger(c)),
		),
	)
	pb.RegisterIdentityServiceServer(
		grpcS,
		server.NewIdentityService(ds, ms, aphgrpc.BaseURLOption(setApiHost(c))),
	)
	reflection.Register(grpcS)

	// http requests muxer
	runtime.HTTPError = aphgrpc.CustomHTTPError
	httpMux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(aphgrpc.HandleCreateResponse),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endP := fmt.Sprintf(":%s", c.String("port"))
	err = pb.RegisterIdentityServiceHandlerFromEndpoint(
		context.Background(),
		httpMux,
		endP,
		opts,
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to register http endpoint for content microservice %s", err),
			2,
		)
	}

	// create listener
	lis, err := net.Listen("tcp", endP)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("failed to listen %s", err),
			2,
		)
	}
	// create the cmux object that will multiplex 2 protocols on same port
	m := cmux.New(lis)
	// match gRPC requests, otherwise regular HTTP requests
	// see https://github.com/grpc/grpc-go/issues/2636#issuecomment-472209287 for why we need to use MatchWithWriters()
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	// CORS setup
	cors := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		OptionsPassthrough: false,
		AllowedHeaders:     []string{"*"},
	})
	httpS := &http.Server{Handler: cors.Handler(httpMux)}
	// collect on this channel the exits of each protocol's .Serve() call
	ech := make(chan error, 2)
	// start the listeners for each protocol
	go func() { ech <- grpcS.Serve(grpcL) }()
	go func() { ech <- httpS.Serve(httpL) }()
	log.Printf("starting multiplexed  server on %s", endP)
	var failed bool
	if err := m.Serve(); err != nil {
		log.Printf("cmux server error: %v", err)
		failed = true
	}
	i := 0
	for err := range ech {
		if err != nil {
			log.Printf("protocol serve error:%v", err)
			failed = true
		}
		i++
		if cap(ech) == i {
			close(ech)
			break
		}
	}
	if failed {
		return cli.NewExitError("error in running cmux server", 2)
	}
	return nil
}

func setApiHost(c *cli.Context) string {
	if len(c.String("identity-api-http-host")) > 0 {
		return c.String("identity-api-http-host")
	}
	return fmt.Sprintf("http://localhost:%s", c.String("port"))
}

func getLogger(c *cli.Context) *logrus.Entry {
	log := logrus.New()
	log.Out = os.Stderr
	switch c.GlobalString("log-format") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	case "json":
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	}
	l := c.GlobalString("log-level")
	switch l {
	case "debug":
		log.Level = logrus.DebugLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	}
	return logrus.NewEntry(log)
}
