// Copyright (c) 2022, Geert JM Vanderkelen

package stores

import "fmt"

type ErrNoObject struct {
	Name string
}

func (e ErrNoObject) Error() string {
	return fmt.Sprintf("%s object not available", e.Name)
}
