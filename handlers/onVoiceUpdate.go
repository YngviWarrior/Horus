package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var voiceStart = map[string]time.Time{}

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

	if vs.ChannelID != "" {
		voiceStart[userID] = time.Now()

		channel, err := s.State.Channel(vs.ChannelID)
		if err != nil {
			channel, err = s.Channel(vs.ChannelID)
			if err != nil {
				fmt.Printf("Erro ao obter o canal: %v\n", err)
				return
			}
		}

		fmt.Printf("[%s] Usuário %s entrou no canal de voz %s\n", time.Now(), member.User.DisplayName(), channel.Name)
	} else {
		start, ok := voiceStart[userID]
		if ok {
			duration := time.Since(start)
			fmt.Printf("[%s] Usuário %s saiu. Ficou por: %s\n", time.Now(), member.User.DisplayName(), duration)
			voiceStart[userID] = time.Time{} // Limpa o tempo de início
		}
	}
}
