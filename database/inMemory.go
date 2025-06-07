package database

import "time"

var VoiceStart = map[string]map[string]time.Time{}
var VoiceTotal = map[string]map[string]time.Duration{}

var PointsTotal = map[string]map[string]int{}

func EnsureMaps(guildID string) {
	if VoiceStart[guildID] == nil {
		VoiceStart[guildID] = map[string]time.Time{}
	}
	if VoiceTotal[guildID] == nil {
		VoiceTotal[guildID] = map[string]time.Duration{}
	}
	if PointsTotal[guildID] == nil {
		PointsTotal[guildID] = map[string]int{}
	}
}
