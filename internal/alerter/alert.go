package alerter

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madjlzz/madprobe/internal/prober"
	"io"
)

// Base type that defines an Alerter.
type Alerter interface {
	Alert(eventBus <-chan prober.Probe)
}

// More specific type of an Alerter that has to close one of it's resource.
type AlertCloser interface {
	Alerter
	io.Closer
}

// Implementation of an alerter that pushes notification to a Discord channel.
type DiscordAlerter struct {
	channelID string
	session   *discordgo.Session
}
