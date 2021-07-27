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

func (d Data) GetObject(key string) *Object {
	for _, o := range d.Objects {
		if o.Key == key {
			return &o
		}
	}
	return nil
}

type Memory struct {
	Data      Data
	Inventory map[string]string
	Solved    []string
}

func createMemory() *Memory {
	var memory Memory
	memory.Data = loadData()
	memory.Inventory = make(map[string]string, 0)
	memory.Solved = make([]string, 0)
	return &memory
}

func (m Memory) ListInventory() {
	if len(m.Inventory) == 0 {
		fmt.Println("Nothing to see here.")
	} else {
		for k, _ := range m.Inventory {
			fmt.Printf("%s\n", k)
		}
	}
}

func (m Memory) IsSolved(object string) bool {
	for _, v := range m.Solved {
		if v == object {
			return true
		}
	}
	return false
}

func (m Memory) GetInventoryItem(item string) string {
	for k, v := range m.Inventory {
		if k == item {
			return v
		}
	}
	return ""
}

func loadData() Data {
	var data Data
	fdata, _ := ioutil.ReadFile("data.json")
	json.Unmarshal(fdata, &data)
	return data
}

func loop(room int, instart bool, memory *Memory) int {
	r := room
	if instart {
		fmt.Println()
		memory.Data.Rooms[room].PrintText()
	}

	fmt.Print("\n> ")

	var cmd string
	var argA string
	var argB string
	var argC string
	fmt.Scanln(&cmd, &argA, &argB, &argC)
	switch cmd {
	case "use":
		if argB == "with" {
			keyWith := memory.GetInventoryItem(argA)
			keyObj := memory.Data.Rooms[room].GetObjectKey(argC)
			if keyObj == "" {
				fmt.Printf("There's no object named %s here.\n", argC)
			} else if keyWith == "" {
				fmt.Printf("You don't have a object named %s in your inventory.\n", argA)
			} else {
				obj := memory.Data.GetObject(keyObj)

				if obj.Command == nil {
					fmt.Printf("%s\n", obj.ActionFailure)
				} else if obj.Command.With != keyWith {
					fmt.Println("I don't think it will work.")
				} else {
					fmt.Printf("%s\n", obj.ActionSuccess)
					if !memory.IsSolved(keyObj) {
						memory.Solved = append(memory.Solved, keyObj)
					}
					if obj.Command.Room != -1 {
						fmt.Print("Press [ENTER] to continue...")
						fmt.Scanln()
						r = obj.Command.Room
					}
				}
			}
		} else if argB == "" {
			keyObj := memory.Data.Rooms[room].GetObjectKey(argA)
			if keyObj == "" {
				fmt.Printf("There's no object named %s here.\n", argA)
			} else {
				obj := memory.Data.GetObject(keyObj)
				if obj.Command == nil || (obj.Command.With != "" && !memory.IsSolved(keyObj)) {
					fmt.Printf("%s\n", obj.ActionFailure)
				} else {
					fmt.Printf("%s\n", obj.ActionSuccess)
					if obj.Command.Room != -1 {
						fmt.Print("Press [ENTER] to continue...")
						fmt.Scanln()
						r = obj.Command.Room
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
			keyObj := memory.Data.Rooms[room].GetObjectKey(argA)
			if keyObj == "" {
				fmt.Printf("There's no object named %s here.\n", argA)
			} else if memory.GetInventoryItem(argA) != "" {
				fmt.Printf("%s already is in inventory.\n", argA)
			} else {
				obj := memory.Data.GetObject(keyObj)
				if obj.IsPocketable {
					memory.Inventory[argA] = keyObj
					fmt.Printf("%s was added to inventory.\n", argA)
				} else {
					fmt.Printf("You can't put %s in your inventory.\n", argA)
				}
			}
		}
	case "check":
		if argB != "" {
			fmt.Println("Invalid command.")
		} else {
			keyObj := memory.Data.Rooms[room].GetObjectKey(argA)
			if keyObj == "" {
				fmt.Printf("There's no object named %s here.\n", argA)
			} else {
				obj := memory.Data.GetObject(keyObj)
				fmt.Printf("%s\n", obj.Description)
			}
		}
	case "inventory":
		memory.ListInventory()
	case "help":
		fmt.Println("Possible commands:")
		fmt.Println("  use OBJECT -> interact with a scene object")
		fmt.Println("  use ITEM with OBJECT -> interact with a scene object using a item from inventory")
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

	return r
}

func Start() {
	memory := createMemory()
	nroom, room, proom := 0, 0, -1
	for {
		nroom = loop(room, proom != room, memory)
		if nroom == -1 {
			break
		} else if nroom == -2 {
			nroom = 0
			memory.Inventory = make(map[string]string)
			memory.Solved = make([]string, 0)
		}

		proom = room
		room = nroom
	}
}
