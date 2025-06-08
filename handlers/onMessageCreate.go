package handlers

import (
	"discord-bot/database"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.GuildID != "" {
		database.PointsTotal[m.GuildID][m.Author.ID] += 1
	}

	if strings.HasPrefix(m.Content, "!help") {
		helpMessage := "Comandos disponíveis:\n" +
			"!help - Mostra esta mensagem de ajuda\n" +
			"!tempo @usuário - Mostra quanto tempo o usuário está em um canal de voz\n" +
			"!invite - Mostra o link de convite do bot\n" +
			"!rank - Mostra o ranking de usuários com mais atividades na semana\n"
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
				startTime, ok := database.VoiceStart[guild.ID][user.ID]
				if !ok {
					if totalTime, exists := database.VoiceTotal[guild.ID][user.ID]; exists {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s não está ativo em um canal de voz, mas tem um total de %s horas nesta semana.", user.DisplayName(), formatDuration(totalTime)))
					} else {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s não está ativo em um canal de voz e não tem tempo registrado nesta semana.", user.DisplayName()))
					}
					return
				}

				duration := time.Since(startTime)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário %s está ativo em voz há %s. Total de %s horas nesta semana.", user.DisplayName(), formatDuration(duration), formatDuration(database.VoiceTotal[m.GuildID][user.ID])))
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

	if strings.HasPrefix(m.Content, "!rank") {
		for _, guild := range s.State.Guilds {
			if guild.ID == m.GuildID {
				if len(m.Mentions) > 0 {
					user := m.Mentions[0]
					points := database.PointsTotal[guild.ID][user.ID]
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuário: %s (%s) - Pontos acumulados: %d", user.Username, user.ID, points))
					return
				}

				pointsMap := database.PointsTotal[m.GuildID]
				var sortedUsers []UserPoints

				for userID, points := range pointsMap {
					sortedUsers = append(sortedUsers, UserPoints{UserID: userID, Points: points})
				}

				// Ordena do maior para o menor
				sort.Slice(sortedUsers, func(i, j int) bool {
					return sortedUsers[i].Points > sortedUsers[j].Points
				})

				// Monta o top 10
				var rankList string = "🏆 **Top 10 Ranking de Pontos**:\n"
				for i, user := range sortedUsers {
					if i >= 10 {
						break
					}

					discordUser, err := s.User(user.UserID)
					if err != nil {
						rankList += fmt.Sprintf("%d. ID %s -> %d pontos\n", i+1, user.UserID, user.Points)
					} else {
						rankList += fmt.Sprintf("%d. %s -> %d pontos\n", i+1, discordUser.DisplayName(), user.Points)
					}
				}

				s.ChannelMessageSend(m.ChannelID, rankList)
				return
			}
		}
	}
}

type UserPoints struct {
	UserID string
	Points int
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
}
