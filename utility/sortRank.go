package utility

import "sort"

type UserPoints struct {
	UserID string
	Points int
}

func SortTotalPoints(pointsMap map[string]int) (sortedUsers []UserPoints) {
	for userID, points := range pointsMap {
		sortedUsers = append(sortedUsers, UserPoints{UserID: userID, Points: points})
	}

	// Ordena do maior para o menor
	sort.Slice(sortedUsers, func(i, j int) bool {
		return sortedUsers[i].Points > sortedUsers[j].Points
	})

	return
}
