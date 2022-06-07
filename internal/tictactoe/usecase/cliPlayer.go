package usecase

type CliPlayer struct {
	score int
}

func (player CliPlayer) GetScore() int {
	return player.score
}

func (player *CliPlayer) SetScore(score int) {
	player.score = score
}
