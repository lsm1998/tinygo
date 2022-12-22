package tinygo

import (
	"github.com/jinzhu/copier"
)

// Deprecated: this function simply calls copier.Copy.
func DeepCopy(dst, src interface{}) error {
	return copier.Copy(dst, src)
}
