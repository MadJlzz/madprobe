package alerter

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/madjlzz/madprobe/internal/prober"
	"log"
	"time"
)

func NewDiscordAlerter() *DiscordAlerter {
	dc, err := NewDiscordConfiguration()
	if err != nil {
		log.Printf("[WARNING] discord configuration contains an error/errors. got: [%v]\n", err)
		return nil
	}
	session, err := discordgo.New("Bot " + dc.Token)
	if err != nil {
		log.Printf("[WARNING] an error occured while trying to initialize session client. got: [%v]\n", err)
		return nil
	}
	err = session.Open()
	if err != nil {
		log.Printf("[WARNING] could not open Websocket to communicate using the session client. got: [%v]\n", err)
		return nil
	}
	return &DiscordAlerter{
		channelID: dc.ChannelID,
		session:   session,
	}
}

func (da *DiscordAlerter) Alert(eventBus <-chan prober.Probe) {
	var probe prober.Probe
	for {
		select {
		case probe = <-eventBus:
			msg := fmt.Sprintf("Probe [%s] is currently [%s]", probe.Name, probe.Status)
			_, err := da.session.ChannelMessageSend(da.channelID, msg)
			if err != nil {
				fmt.Println(err)
			}
		default:
			time.Sleep(time.Duration(probe.Delay) * time.Second)
		}
	}
}

func (da *DiscordAlerter) Close() error {
	return da.session.Close()
}
