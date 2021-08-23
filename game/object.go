package game

type Command struct {
	With      string
	Room      int
	Condition string
}

type Object struct {
	Key           string
	Name          string
	Description   string
	ActionSuccess string
	ActionFailure string
	IsPocketable  bool
	Commands      []Command
}

func (o Object) GetCommandsWith(with string) []Command {
	var ans []Command
	ans = make([]Command, 0)
	for _, v := range o.Commands {
		if v.With == with {
			ans = append(ans, v)
		}
	}
	return ans
}
