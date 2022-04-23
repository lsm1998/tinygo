package tinygo

import (
	"github.com/google/uuid"
	"strings"
)

func UUID() string {
	return uuid.New().String()
}

func IsBlank(str string) bool {
	return len(strings.ReplaceAll(str, " ", "")) == 0
}
