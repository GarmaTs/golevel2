package common

type IUserInteract interface {
	GetUserMove(bool) (int, int)
	SetFieldSize() int
	ShowError(err error)
	ShowField(size int, field [][]string)
	ShowGameResult(res string)
	AskForRepeat() bool
	ShowTotalScore(res string)
}
