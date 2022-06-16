package common

type IUserInteract interface {
	GetUserMove(bool, bool) (int, int)
	SetFieldSize() int
	ShowError(err error)
	ShowField(size int, field [][]string)
	ShowGameResult(res string)
	AskForRepeat() bool
	ShowTotalScore(res string)
	IsVersusAI() bool
	ShowMessage(msg string)
}
