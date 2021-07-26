package game

type Command struct {
	With string
	Room int
}

type Object struct {
	Name          string
	Description   string
	ActionSuccess string
	ActionFailure string
	IsPocketable  bool
	Commands      []Command
}
