package alerter

import (
	"errors"
	"github.com/spf13/viper"
)

var ErrDiscordChannelNotValid = errors.New("channel id must be set")
var ErrDiscordTokenNotValid = errors.New("token must be set")

// Discord struct holding default configuration option.
type DiscordConfiguration struct {
	// The ID of the Discord channel used for posting new alerts.
	ChannelID string
	// The authentication Token to talk with the Discord API.
	Token string
}

// Default value of the ServerConfiguration struct.
var DefaultDiscordConfiguration = &DiscordConfiguration{
	ChannelID: "",
	Token:     "",
}

func NewDiscordConfiguration() (*DiscordConfiguration, error) {
	dc := &DiscordConfiguration{
		ChannelID: viper.GetString("discord-channel-id"),
		Token:     viper.GetString("discord-token"),
	}
	return dc, dc.validate()
}

func (dc *DiscordConfiguration) validate() error {
	if len(dc.ChannelID) <= 0 {
		return ErrDiscordChannelNotValid
	}
	if len(dc.Token) <= 0 {
		return ErrDiscordTokenNotValid
	}
	return nil
}
