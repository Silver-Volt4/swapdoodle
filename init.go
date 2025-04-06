package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/PretendoNetwork/grpc/go/account"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"

	"github.com/PretendoNetwork/nex-go/v2"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func init() {
	globals.Logger = plogger.NewLogger()

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	s3Endpoint := os.Getenv("PN_SD_CONFIG_S3_ENDPOINT")
	s3AccessKey := os.Getenv("PN_SD_CONFIG_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv("PN_SD_CONFIG_S3_ACCESS_SECRET")
	s3SecureEnv := os.Getenv("PN_SD_CONFIG_S3_SECURE")

	postgresURI := os.Getenv("PN_SD_POSTGRES_URI")
	hppServerPort := os.Getenv("PN_SD_HPP_SERVER_PORT")
	accountGRPCHost := os.Getenv("PN_SD_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv("PN_SD_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv("PN_SD_ACCOUNT_GRPC_API_KEY")

	if strings.TrimSpace(postgresURI) == "" {
		globals.Logger.Error("PN_SD_POSTGRES_URI environment variable not set")
		os.Exit(0)
	}

	kerberosPassword := make([]byte, 0x10)
	_, err = rand.Read(kerberosPassword)
	if err != nil {
		globals.Logger.Error("Error generating Kerberos password")
		os.Exit(0)
	}

	globals.KerberosPassword = string(kerberosPassword)

	globals.HppServerAccount = nex.NewAccount(2, "Quazal Rendez-Vous", globals.KerberosPassword)

	if strings.TrimSpace(hppServerPort) == "" {
		globals.Logger.Error("PN_SD_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(hppServerPort); err != nil {
		globals.Logger.Errorf("PN_SD_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", hppServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_SD_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", hppServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCHost) == "" {
		globals.Logger.Error("PN_SD_ACCOUNT_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCPort) == "" {
		globals.Logger.Error("PN_SD_ACCOUNT_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(accountGRPCPort); err != nil {
		globals.Logger.Errorf("PN_SD_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_SD_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_SD_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCAccountClientConnection, err = grpc.NewClient(fmt.Sprintf("%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)

	staticCredentials := credentials.NewStaticV4(s3AccessKey, s3AccessSecret, "")

	s3Secure, err := strconv.ParseBool(s3SecureEnv)
	if err != nil {
		globals.Logger.Warningf("PN_SD_CONFIG_S3_SECURE environment variable not set. Using default value: %t", true)
		s3Secure = true
	}

	minIOClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  staticCredentials,
		Secure: s3Secure,
	})
	if err != nil {
		panic(err)
	}

	globals.MinIOClient = minIOClient
	globals.Presigner = globals.NewS3Presigner(globals.MinIOClient)

	// * Connect to and setup databases
	database.ConnectPostgres()
}
