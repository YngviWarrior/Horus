package jobs

import (
	"discord-bot/database"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func InitPointsTotal(s *discordgo.Session, guildID string) {
	after := ""
	limit := 1000 // Máximo permitido por chamada

	for {
		members, err := s.GuildMembers(guildID, after, limit)
		if err != nil {
			fmt.Printf("Erro ao buscar membros: %v\n", err)
			break
		}

		if len(members) == 0 {
			break
		}

		for _, member := range members {
			userID := member.User.ID
			after = userID
			database.PointsTotal[guildID][userID] = 0
		}

		// Discord retorna em blocos — se menos de 1000, acabou
		if len(members) < limit {
			break
		}
	}
	fmt.Println("População de usuários para pontos total concluída para o servidor:", guildID)
}
