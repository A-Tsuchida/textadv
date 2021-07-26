package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Data struct {
	Objects []Object
	Rooms   []Room
}

func loadData() Data {
	var data Data
	fdata, _ := ioutil.ReadFile("data.json")
	json.Unmarshal(fdata, &data)
	return data
}

func loop(room int, instart bool, data *Data, inventory []string) (int, []string) {
	r := room
	if instart {
		fmt.Println()
		data.Rooms[room].PrintText()
	}

	fmt.Print("\n> ")

	var cmd string
	var argA string
	var argB string
	var argC string
	fmt.Scanln(&cmd, &argA, &argB, &argC)
	switch cmd {
	case "use":
		finished := false
		err := true
		if argB == "with" {
			ininv := false
			for _, v2 := range inventory {
				if v2 == argA {
					ininv = true
				}
			}
			if !ininv {
				fmt.Printf("You don't have a object named %s in your inventory.", argA)
				err = false
			} else {
				for k, v := range data.Rooms[room].Objects {
					if argC == k {
						obj := data.Objects[v]
						for _, c := range obj.Commands {
							if c.With == argA {
								fmt.Printf("%s\n", obj.ActionSuccess)
								if c.Room != -1 {
									fmt.Print("Press [ENTER] to continue...")
									fmt.Scanln()
									r = c.Room
								}
								finished = true
								err = false
							}
						}
						if !finished {
							fmt.Printf("%s\n", obj.ActionFailure)
							err = false
						}
					}
				}
			}
			if err {
				fmt.Println("Invalid command.")
			}
		} else if argB == "" {
			for k, v := range data.Rooms[room].Objects {
				if argA == k {
					obj := data.Objects[v]
					for _, c := range obj.Commands {
						if c.With == "" {
							fmt.Printf("%s\n", obj.ActionSuccess)
							if c.Room != -1 {
								fmt.Print("Press [ENTER] to continue...")
								fmt.Scanln()
								r = c.Room
							}
							finished = true
							err = false
						}
					}
					if !finished {
						fmt.Printf("%s\n", obj.ActionFailure)
						err = false
					}
				}
			}
		} else {
			fmt.Println("Invalid command.")
		}
	case "get":
		if argB != "" {
			fmt.Println("Invalid command.")
		} else {
			for k, v := range data.Rooms[room].Objects {
				obj := data.Objects[v]
				if argA == k {
					if obj.IsPocketable {
						found := false
						for _, v := range inventory {
							if v == argA {
								found = true
							}
						}
						if !found {
							fmt.Printf("%s was added to inventory.\n", obj.Name)
							inventory = append(inventory, obj.Name)
						} else {
							fmt.Printf("%s already is in inventory.\n", obj.Name)
						}
					} else {
						fmt.Println("Invalid command for this object.")
					}
				}
			}
		}
	case "check":
		if argB != "" {
			fmt.Println("Invalid command.")
		} else {
			for k, v := range data.Rooms[room].Objects {
				obj := data.Objects[v]
				if argA == k {
					fmt.Printf("%s\n", obj.Description)
				}
			}
		}
	case "inventory":
		if len(inventory) == 0 {
			fmt.Println("Nothing to see here.")
		} else {
			for _, v := range inventory {
				fmt.Printf("%s\n", v)
			}
		}
	case "help":
		fmt.Println("Possible commands:")
		fmt.Println("  use OBJECT -> interact with a scene object")
		fmt.Println("  get OBJECT -> puts OBJECT in your inventory")
		fmt.Println("  check OBJECT -> get OBJECT description")
		fmt.Println("  inventory -> list what is in your inventory")
		fmt.Println("  help -> command list")
		fmt.Println("  save FILENAME -> save current game to FILENAME")
		fmt.Println("  load FILENAME -> load game stored in FILENAME")
		fmt.Println("  newgame -> start a new game")
	case "save":
		fmt.Fprintln(os.Stderr, "Not Implemented.")
	case "load":
		fmt.Fprintln(os.Stderr, "Not Implemented.")
	case "newgame":
		fmt.Print("Are you sure? (Y/N) ")
		var ans string
		fmt.Scanln(&ans)
		if ans == "Y" || ans == "y" {
			r = -2
		}
	case "quit":
		fmt.Print("Are you sure? (Y/N) ")
		var ans string
		fmt.Scan(&ans)
		if ans == "Y" || ans == "y" {
			r = -1
		}

	default:
		fmt.Printf("There in no command %s. Type 'help' for command list.\n", cmd)
	}

	return r, inventory
}

func Start() {
	data := loadData()
	inventory := make([]string, 0)
	nroom, room, proom := 0, 0, -1
	for {
		nroom, inventory = loop(room, proom != room, &data, inventory)
		if nroom == -1 {
			break
		} else if nroom == -2 {
			nroom = 0
			inventory = make([]string, 0)
		}

		proom = room
		room = nroom
	}
}
