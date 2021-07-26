package game

import (
	"fmt"
)

type Room struct {
	Name    string
	Text    string
	Objects map[string]int
}

func (r Room) PrintText() {
	fmt.Printf("%s\n%s\n", r.Name, r.Text)
}
