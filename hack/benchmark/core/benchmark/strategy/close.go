package strategy

type CloseOption func(*close)

var (
	defaultCloseOptions = []CloseOption{}
)
