package util

import (
	"github.com/spf13/viper"
	"time"
)

type ServerConfiguration struct {
	// the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m
	Wait time.Duration
	// the port for which the server will start to listen to
	Port string
	// public certificate shown by the server to it's clients
	ServerCertificate string
	// the server's certificate private key
	ServerKey string
	// the CA certificate
	CaCertificate string
}

// Default value of the ServerConfiguration struct.
var DefaultServerConfiguration = &ServerConfiguration{
	Wait:              time.Second * 15,
	Port:              "3000",
	ServerCertificate: "",
	ServerKey:         "",
	CaCertificate:     "",
}

// Insert a new ServerConfiguration with default values or values coming from Viper.
func NewServerConfiguration() *ServerConfiguration {
	return &ServerConfiguration{
		Wait:              viper.GetDuration("graceful-timeout"),
		Port:              viper.GetString("port"),
		ServerCertificate: viper.GetString("cert"),
		ServerKey:         viper.GetString("key"),
		CaCertificate:     viper.GetString("ca-cert"),
	}
}
