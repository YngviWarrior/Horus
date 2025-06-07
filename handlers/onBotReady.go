package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func OnBotReady(s *discordgo.Session, m *discordgo.Ready) {
	clientID := s.State.User.ID // ID do próprio bot

	const permissions = 8 // admin

	inviteURL := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&permissions=%d&scope=bot+applications.commands", clientID, permissions)

	fmt.Printf("Adicione o bot em seu servidor usando este link:\n%s", inviteURL)
}
