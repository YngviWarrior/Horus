package handlers

import (
	"discord-bot/database"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnVoiceUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	userID := vs.Member.User.ID
	guildID := vs.GuildID

	member, err := s.State.Member(guildID, userID)
	if err != nil {
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			fmt.Printf("Erro ao obter membro: %v\n", err)
			return
		}
	}

	if vs.SelfDeaf || vs.SelfMute || vs.Mute || vs.Deaf {
		start, ok := database.VoiceStart[guildID][userID]
		if ok && !start.IsZero() {
			duration := time.Since(start)
			database.VoiceTotal[guildID][userID] += duration

			fmt.Printf("[%s] Usuário %s Mutou. Sessão: %s | Total acumulado no servidor: %s\n",
				time.Now().Format(time.RFC3339),
				member.User.Username,
				duration,
				database.VoiceTotal[guildID][userID],
			)
			delete(database.VoiceStart[guildID], userID)
		}
	} else if vs.ChannelID != "" {
		database.VoiceStart[guildID][userID] = time.Now()

		channel, err := s.State.Channel(vs.ChannelID)
		if err != nil {
			channel, err = s.Channel(vs.ChannelID)
			if err != nil {
				fmt.Printf("Erro ao obter o canal: %v\n", err)
				return
			}
		}

		guild, err := s.State.Guild(guildID)
		if err != nil {
			guild, err = s.Guild(guildID)
			if err != nil {
				fmt.Printf("Erro ao obter o servidor: %v\n", err)
				return
			}
		}

		fmt.Printf("[%s] Usuário %s entrou no canal de voz %s no server %s\n", time.Now(), member.User.DisplayName(), channel.Name, guild.Name)
	} else {
		start, ok := database.VoiceStart[guildID][userID]
		if ok && !start.IsZero() {
			duration := time.Since(start)
			database.VoiceTotal[guildID][userID] += duration

			fmt.Printf("[%s] Usuário %s saiu. Sessão: %s | Total acumulado no servidor: %s\n",
				time.Now().Format(time.RFC3339),
				member.User.Username,
				duration,
				database.VoiceTotal[guildID][userID],
			)
			delete(database.VoiceStart[guildID], userID)
		}
	}

}
