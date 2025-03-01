package service

type History interface {
	Write(obj any)
}

type HistoryDummy struct{}

func (h HistoryDummy) Write(any) {}
