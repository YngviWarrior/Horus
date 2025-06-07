package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "!help") {
		helpMessage := "Comandos disponíveis:\n" +
			"!help - Mostra esta mensagem de ajuda\n" +
			"!tempo @usuário - Mostra quanto tempo o usuário está em um canal de voz\n" +
			"!invite - Mostra o link de convite do bot\n"
		s.ChannelMessageSend(m.ChannelID, helpMessage)
		return
	}

	// Comando "!tempo @usuário"
	if strings.HasPrefix(m.Content, "!tempo") {
		for _, guild := range s.State.Guilds {
			if guild.ID == m.GuildID {
				if len(m.Mentions) == 0 {
					s.ChannelMessageSend(m.ChannelID, "Por favor, mencione um usuário para verificar o tempo.")
					return
				}

				user := m.Mentions[0]
				startTime, ok := voiceStart[guild.ID][user.ID]
				if !ok {
					if totalTime, exists := voiceTotal[guild.ID][user.ID]; exists {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s não está em um canal de voz, mas tem um total de %s horas nesta semana.", user.Username, formatDuration(totalTime)))
					} else {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s não está em um canal de voz e não tem tempo registrado nesta semana.", user.Username))
					}
					return
				}

				duration := time.Since(startTime)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s está em voz há %s. Total de %s horas nesta semana.", user.DisplayName(), formatDuration(duration), formatDuration(voiceTotal[m.GuildID][user.ID])))
			}
		}
	}

	// Comando "!invite"
	if strings.HasPrefix(m.Content, "!invite") {
		clientID := s.State.User.ID // ID do próprio bot
		const permissions = 8       // admin

		inviteURL := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&permissions=%d&scope=bot+applications.commands", clientID, permissions)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Adicione o bot em seu servidor usando este link:\n%s", inviteURL))
	}
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
}
