package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/proto/sqlQuery"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"time"
)

type server struct {
	sqlQuery.UnimplementedQueryRunnerServer
}

func createErrorReply(err error, encryptPath string) (*sqlQuery.SqlRunnerReply, error) {
	errStr := "Failed to run query: " + err.Error()

	buffer := []byte(errStr)
	if encryptPath != "" {
		var newErr error
		buffer, newErr = vault.TransitEncrypt(encryptPath, buffer)
		if newErr != nil {
			return createErrorReply(newErr, "")
		}
	}

	core.Info("\tFailure Reply: ", buffer)
	return &sqlQuery.SqlRunnerReply{
		EncryptedData: buffer,
		Success:       false,
	}, nil
}

func (s *server) RunSqlQuery(ctx context.Context, in *sqlQuery.SqlRunnerRequest) (*sqlQuery.SqlRunnerReply, error) {
	core.Info("Run Query: ", in.QueryId, in.OrgId, in.VaultResultPath)
	// Ignore this error. Should be fine?
	err := vault.TransitCreateNewEngineKey(in.VaultResultPath)
	if err != nil {
		return createErrorReply(err, "")
	}

	result, err := runQuery(in.QueryId, in.OrgId)
	if err != nil {
		return createErrorReply(err, in.VaultResultPath)
	}

	marshal, err := json.Marshal(result)
	if err != nil {
		return createErrorReply(err, in.VaultResultPath)
	}

	encrypted, err := vault.TransitEncrypt(in.VaultResultPath, marshal)
	if err != nil {
		return createErrorReply(err, "")
	}

	core.Info("\tSuccess Reply: ", encrypted)
	return &sqlQuery.SqlRunnerReply{
		EncryptedData: encrypted,
		Success:       true,
	}, nil
}

func main() {
	core.Init()
	database.Init()
	webcore.InitializeWebcore()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())

	queryId := flag.Int64("queryId", -1, "Refresh ID to retrieve data for. Will not read from gRPC if specified.")
	orgId := flag.Int64("orgId", -1, "Org ID to retrieve data for. Will not read from gRPC if specified.")
	rpcClient := flag.Bool("rpc", false, "If query ID and org ID specified, will use an RPC client to execute.")
	flag.Parse()

	if *queryId >= 0 && *orgId >= 0 {
		if *rpcClient {
			conn, err := webcore.CreateGRPCClientConnection(core.EnvConfig.Grpc.QueryRunner, core.EnvConfig.Tls)
			if err != nil {
				core.Error("Failed to connect to GRPC: " + err.Error())
			}

			defer conn.Close()

			client := sqlQuery.NewQueryRunnerClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			transitKey := fmt.Sprintf("sqlquery-%d", *queryId)

			resp, err := client.RunSqlQuery(ctx, &sqlQuery.SqlRunnerRequest{
				VaultResultPath: transitKey,
				QueryId:         *queryId,
				OrgId:           int32(*orgId),
			})

			if err != nil {
				core.Error("Failed to run query: " + err.Error())
			}

			decrypt, err := vault.TransitDecrypt(transitKey, resp.EncryptedData)
			if err != nil {
				core.Error("Failed to decrypt result: " + err.Error())
			}

			core.Info(string(decrypt))
			core.Info(resp.Success)
		} else {
			result, err := runQuery(*queryId, int32(*orgId))
			if err != nil {
				core.Error("Failed to run query: " + err.Error())
			}
			core.Info(result.Columns)
			core.Info(result.CsvText)
		}
	} else {
		lis, s, err := webcore.CreateGRPCServer(core.EnvConfig.Grpc.QueryRunner)
		if err != nil {
			core.Error("Failed to create GRPC server: " + err.Error())
		}

		sqlQuery.RegisterQueryRunnerServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			core.Error("Failed to serve: " + err.Error())
		}
	}
}
