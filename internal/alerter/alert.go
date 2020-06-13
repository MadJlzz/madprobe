package alerter

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madjlzz/madprobe/internal/service"
)

type Alerter interface {
	Alert(eventBus <-chan service.Probe)
}

type DiscordAlerter struct {
	channelID string
	session   *discordgo.Session
}
