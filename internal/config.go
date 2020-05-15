package internal

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"time"
)

const (
	defaultWaitDuration      = time.Second * 15
	defaultPort              = "3000"
	defaultServerCertificate = ""
	defaultServerKey         = ""
	defaultCaCertificate     = ""
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

func NewServerConfiguration() ServerConfiguration {
	initConfigurationFiles()
	initEnvironmentVariables()
	initFlags()
	return ServerConfiguration{
		Wait:              getWaitTime(),
		Port:              getPort(),
		ServerCertificate: getServerCertificate(),
		ServerKey:         getServerKey(),
		CaCertificate:     getCaCertificate(),
	}
}

func initConfigurationFiles() {
	viper.SetConfigName("madprobe")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[WARNING] %v\n", err)
	}
}

func initEnvironmentVariables() {
	viper.AutomaticEnv()
}

func initFlags() {
	flag.Duration("graceful-timeout", defaultWaitDuration, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.String("port", defaultPort, "the port for which the server will start to listen to")
	flag.String("cert", defaultServerCertificate, "public certificate shown by the server to it's clients")
	flag.String("key", defaultServerKey, "the server's certificate private key")
	flag.String("ca-cert", defaultCaCertificate, "the CA certificate")
	flag.Parse()
	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		log.Printf("[WARNING] %v\n", err)
	}
}

func getWaitTime() time.Duration {
	waitTime := defaultWaitDuration
	if viper.IsSet("graceful-timeout") {
		waitTime = viper.GetDuration("graceful-timeout")
	}
	return waitTime
}

func getPort() string {
	port := defaultPort
	if viper.IsSet("port") {
		port = viper.GetString("port")
	}
	return port
}

func getServerCertificate() string {
	serverCertificate := defaultServerCertificate
	if viper.IsSet("cert") {
		serverCertificate = viper.GetString("cert")
	}
	return serverCertificate
}

func getServerKey() string {
	serverKey := defaultServerKey
	if viper.IsSet("key") {
		serverKey = viper.GetString("key")
	}
	return serverKey
}

func getCaCertificate() string {
	caCert := defaultCaCertificate
	if viper.IsSet("ca-cert") {
		caCert = viper.GetString("ca-cert")
	}
	return caCert
}
