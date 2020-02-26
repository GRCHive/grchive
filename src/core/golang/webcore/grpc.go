package webcore

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

func CreateGRPCClientConnection(cfg core.GrpcConfig, tls *core.TLSConfig, extra ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts := append([]grpc.DialOption{grpc.WithBlock()}, extra...)
	if cfg.TLS {
		creds := credentials.NewTLS(tls.Config())
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	url := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	core.Info(url, cfg.TLS)
	return grpc.Dial(url, opts...)
}

func CreateGRPCServer(cfg core.GrpcConfig) (net.Listener, *grpc.Server, error) {
	url := fmt.Sprintf(":%d", core.EnvConfig.Grpc.QueryRunner.Port)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return nil, nil, err
	}

	opts := []grpc.ServerOption{}
	if cfg.TLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCert, cfg.TLSKey)
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	core.Info(url, cfg.TLS)
	return lis, grpc.NewServer(opts...), nil
}
