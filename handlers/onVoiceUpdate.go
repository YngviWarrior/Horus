package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var voiceStart = map[string]time.Time{}
var voiceTotal = map[string]time.Duration{}

func ResetVoiceData(timeAlive time.Time) {
	if time.Since(timeAlive) < 24*time.Hour*7 {
		return // Não reseta se o bot está ativo há menos de 7 dias
	}

	voiceStart = make(map[string]time.Time)
	voiceTotal = make(map[string]time.Duration)
}

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
		if ok && !start.IsZero() {
			duration := time.Since(start)
			voiceTotal[userID] += duration // Acumula tempo
			fmt.Printf("[%s] Usuário %s saiu. Sessão: %s | Total: %s\n",
				time.Now().Format(time.RFC3339),
				member.User.Username,
				duration,
				voiceTotal[userID],
			)
			delete(voiceStart, userID)
		}
	}
}
