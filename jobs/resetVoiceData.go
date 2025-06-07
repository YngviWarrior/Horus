package jobs

import (
	"discord-bot/database"
	"time"
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

func ResetVoiceData(timeAlive time.Time) {
	ticker := time.NewTicker(7 * 24 * time.Hour) // 7 dias
	defer ticker.Stop()
	for {
		<-ticker.C
		database.VoiceStart = map[string]map[string]time.Time{}
		database.VoiceTotal = map[string]map[string]time.Duration{}
	}
}
