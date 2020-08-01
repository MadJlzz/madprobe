package util

import (
	"github.com/madjlzz/madprobe/internal/alerter"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

// Viper flag set used for the default application args.
var ViperFlagSet = flag.CommandLine

func init() {
	serverFlags()
	discordFlags()
	parse()
}

func serverFlags() {
	ViperFlagSet.Duration("graceful-timeout", DefaultServerConfiguration.Wait, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	ViperFlagSet.String("port", DefaultServerConfiguration.Port, "the port for which the server will start to listen to")
	ViperFlagSet.String("cert", DefaultServerConfiguration.ServerCertificate, "public certificate shown by the server to it's clients")
	ViperFlagSet.String("key", DefaultServerConfiguration.ServerKey, "the server's certificate private key")
	ViperFlagSet.String("ca-cert", DefaultServerConfiguration.CaCertificate, "the CA certificate")
}

func discordFlags() {
	ViperFlagSet.String("discord-channel-id", alerter.DefaultDiscordConfiguration.ChannelID, "the Discord channel for posting alerts")
	ViperFlagSet.String("discord-token", alerter.DefaultDiscordConfiguration.Token, "the Discord Token for authentication")
}

func parse() {
	flag.Parse()
	if err := viper.BindPFlags(ViperFlagSet); err != nil {
		log.Printf("[WARNING] an error occured while binding application parse. got: [%v]\n", err)
	}
}
