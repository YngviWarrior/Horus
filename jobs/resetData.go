package jobs

import (
	"discord-bot/database"
	"discord-bot/utility"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func AddVoicePoints(guildID string) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		for userID := range database.VoiceStart[guildID] {
			database.PointsTotal[guildID][userID] += 10
		}
	}
}

func FindFirstTextChannel(s *discordgo.Session, guildID string) string {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		log.Println("Erro ao obter canais:", err)
		return ""
	}

	for _, c := range channels {
		if c.Type == discordgo.ChannelTypeGuildText && c.ID != "" {
			return c.ID
		}
	}

	return ""
}

func ResetData(s *discordgo.Session, timeAlive time.Time) {
	ticker := time.NewTicker(7 * (24 * time.Hour)) // 7 dias
	defer ticker.Stop()

	for {
		<-ticker.C

		for _, guild := range s.State.Guilds {
			channelID := FindFirstTextChannel(s, guild.ID)
			if channelID != "" {

				pointsMap := database.PointsTotal[guild.ID]
				sortedUsers := utility.SortTotalPoints(pointsMap)

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

				s.ChannelMessageSend(channelID, rankList)

				_, err := s.ChannelMessageSend(channelID, "🔄 Dados de voz estão sendo resetados agora neste servidor.")
				if err != nil {
					log.Println("Erro ao enviar mensagem no canal:", err)
				}
			}
		}

		// 🔥 Faz o reset
		database.VoiceStart = map[string]map[string]time.Time{}
		database.VoiceTotal = map[string]map[string]time.Duration{}
		database.PointsTotal = map[string]map[string]int{}

		log.Println("Dados de voz resetados.")
	}
}
