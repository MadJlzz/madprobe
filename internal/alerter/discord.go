package alerter

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/madjlzz/madprobe/internal/prober"
	"github.com/spf13/viper"
	"log"
	"time"
)

func NewDiscordAlerter() *DiscordAlerter {
	channelId := viper.GetString("discord-channel-id")
	token := viper.GetString("discord-token")

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("an unexpected error occured while trying to initialize discord client")
		return nil
	}

	err = discord.Open()
	if err != nil {
		log.Println(err)
	}

	return &DiscordAlerter{
		channelID: channelId,
		session:   discord,
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