package game

import (
	"fmt"
)

type Room struct {
	Name    string
	Text    string
	Objects map[string]string
}

func (r Room) PrintText() {
	fmt.Printf("\x1b[1m%s\x1b[0m\n%s\n", r.Name, r.Text)
}

func (r Room) GetObjectKey(o string) string {
	for k, v := range r.Objects {
		if o == k {
			return v
		}
	}
	return ""
}
