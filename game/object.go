package game

type Command struct {
	With string
	Room int
}

type Object struct {
	Key           string
	Name          string
	Description   string
	ActionSuccess string
	ActionFailure string
	IsPocketable  bool
	Command       *Command
}
