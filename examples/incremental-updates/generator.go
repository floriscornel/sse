package main

import (
	"fmt"
	"math/rand"
)

func generatePlayerList(size int) map[int]playerScore {
	players := make(map[int]playerScore, size)
	for i := 0; i < size; i++ {
		player := generateRandomPlayer()
		players[player.ID] = player
	}
	return players
}

func playerSliceFromMap(m map[int]playerScore) []playerScore {
	var players []playerScore
	for _, v := range m {
		players = append(players, v)
	}
	return players
}

func playerIDSliceFromMap(m map[int]playerScore) []int {
	var ids []int
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}

func generateRandomPlayer() playerScore {
	return playerScore{
		ID:    generateRandomPlayerID(),
		Name:  generateRandomPlayerName(),
		Score: generateRandomPlayerScore(),
	}
}

func generateRandomFirstName() string {
	firstNames := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi",
		"Ivan", "Judy", "Kevin", "Linda", "Mallory", "Nancy", "Oscar", "Peggy", "Quentin",
		"Romeo", "Sierra", "Tango", "Ursula", "Victor", "Wendy", "Xander", "Yvonne", "Zelda"}

	return firstNames[rand.Intn(len(firstNames))]
}

func generateRandomLastName() string {
	lastNames := []string{"Adams", "Brown", "Clark", "Davis", "Evans", "Frank", "Ghosh", "Hills",
		"Iyer", "Jones", "Klein", "Lopez", "Mason", "Nguyen", "Owens", "Patel", "Quinn",
		"Reed", "Smith", "Taylor", "Unger", "Vargas", "Wong", "Xu", "Yilmaz", "Zhang"}

	return lastNames[rand.Intn(len(lastNames))]
}

func generateRandomPlayerName() string {
	return fmt.Sprintf("%s %s", generateRandomFirstName(), generateRandomLastName())
}

func generateRandomPlayerScore() int {
	return rand.Intn(100)
}

func generateRandomPlayerID() int {
	return rand.Intn(999)
}
