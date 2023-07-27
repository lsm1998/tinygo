package logx

import (
	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

func NewEntry() *Entry {
	return &Entry{logrus.NewEntry(logrus.StandardLogger())}
}
