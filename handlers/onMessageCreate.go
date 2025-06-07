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

	// Comando "!tempo @usuário"
	if strings.HasPrefix(m.Content, "!tempo") {
		if len(m.Mentions) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Por favor, mencione um usuário para verificar o tempo.")
			return
		}

		user := m.Mentions[0]
		startTime, ok := voiceStart[user.ID]
		if !ok {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s não está em um canal de voz.", user.Username))
			return
		}

		duration := time.Since(startTime)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s está em voz há %s.", user.DisplayName(), formatDuration(duration)))
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
