package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnBotReady(s *discordgo.Session, m *discordgo.Ready) {
	clientID := s.State.User.ID // ID do próprio bot

	const permissions = 8 // admin

	inviteURL := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&permissions=%d&scope=bot+applications.commands", clientID, permissions)

	fmt.Printf("Adicione o bot em seu servidor usando este link:\n%s", inviteURL)

	for _, guild := range s.State.Guilds {
		channels, err := s.GuildChannels(guild.ID)
		if err != nil {
			fmt.Printf("Erro ao obter canais da guild %s: %v\n", guild.ID, err)
			continue
		}

		ensureGuildMaps(guild.ID)

		for _, server := range channels {
			if server.Type == discordgo.ChannelTypeGuildVoice {
				for _, vs := range guild.VoiceStates {
					if vs.ChannelID == server.ID {
						voiceStart[guild.ID][vs.UserID] = time.Now()
						user, err := s.User(vs.UserID)
						if err == nil {
							fmt.Printf("Usuário %s (%s) já está no canal de voz %s no server %s\n", user.Username, user.ID, server.Name, guild.Name)
						} else {
							fmt.Printf("Usuário %s já está no canal de voz %s no server %s\n", vs.UserID, server.Name, guild.Name)
						}
					}
				}
			}
		}
	}
}
