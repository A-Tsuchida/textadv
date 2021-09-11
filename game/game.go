package game

import (
	"fmt"
)

func createMemory() *Memory {
	var memory Memory
	memory.Assets = LoadAssets()
	memory.User.Inventory = make(map[string]string, 0)
	memory.User.Solved = make([]string, 0)
	return &memory
}

func printError(message string) {
	fmt.Println("\x1b[1;31m" + message + "\x1b[0m")
}

func printWarning(message string) {
	fmt.Println("\x1b[1;33m" + message + "\x1b[0m")
}

func printOK(message string) {
	fmt.Println("\x1b[1;32m" + message + "\x1b[0m")
}

func printNextWait() {
	printOK("Press [ENTER] to continue...")
	fmt.Scanln()
}

func loop(instart bool, memory *Memory) int {
	r := memory.User.Room
	if instart {
		fmt.Print("\x1b[2J\x1b[H")
		fmt.Println()
		memory.Assets.Rooms[memory.User.Room].PrintText()
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
			keyObj := memory.Assets.Rooms[memory.User.Room].GetObjectKey(argA)
			keyWith := memory.User.GetInventoryItem(argC)
			if keyObj == "" {
				printWarning(fmt.Sprintf("There's no object named %s here.\n", argA))
			} else if keyWith == "" {
				printWarning(fmt.Sprintf("You don't have a object named %s in your inventory.\n", argC))
			} else {
				obj := memory.Assets.GetObject(keyObj)

				if len(obj.Commands) == 0 {
					printWarning(obj.ActionFailure)
				} else {
					commands := obj.GetCommandsWith(keyWith)
					if len(commands) == 0 {
						printWarning("I don't think it will work.")
					} else {
						for _, c := range commands {
							if memory.User.IsSolved(c.Condition) {
								printOK(obj.ActionSuccess)
								if !memory.User.IsSolved(keyObj) {
									memory.User.Solved = append(memory.User.Solved, keyObj)
								}
								if c.Room != -1 {
									printNextWait()
									r = c.Room
								}
							}
						}
					}
				}
			}
		} else if argB == "" {
			keyObj := memory.Assets.Rooms[memory.User.Room].GetObjectKey(argA)
			if keyObj == "" {
				printWarning(fmt.Sprintf("There's no object named %s here.\n", argA))
			} else {
				obj := memory.Assets.GetObject(keyObj)
				if len(obj.Commands) > 0 {
					commands := obj.GetCommandsWith("")
					if len(commands) > 0 {
						solved := false
						for _, c := range commands {
							if memory.User.IsSolved(c.Condition) {
								printOK(obj.ActionSuccess)
								if !memory.User.IsSolved(keyObj) {
									memory.User.Solved = append(memory.User.Solved, keyObj)
								}
								if c.Room != -1 {
									printNextWait()
									r = c.Room
								}
								solved = true
								break
							}
						}
						if !solved {
							printWarning(obj.ActionFailure)
						}
					} else {
						printWarning(obj.ActionFailure)
					}
				} else {
					printWarning(obj.ActionFailure)
				}
			}
		} else {
			printError("Invalid command.")
		}
	case "get":
		if argB != "" {
			printError("Invalid command.")
		} else {
			keyObj := memory.Assets.Rooms[memory.User.Room].GetObjectKey(argA)
			if keyObj == "" {
				printWarning(fmt.Sprintf("There's no object named %s here.", argA))
			} else if memory.User.GetInventoryItem(argA) != "" {
				printWarning(fmt.Sprintf("%s already is in inventory.", argA))
			} else {
				obj := memory.Assets.GetObject(keyObj)
				if obj.IsPocketable {
					memory.User.Inventory[argA] = keyObj
					printOK(fmt.Sprintf("%s was added to inventory.", argA))
				} else {
					printWarning(fmt.Sprintf("You can't put %s in your inventory.", argA))
				}
			}
		}
	case "check":
		if argB != "" {
			printError("Invalid command.")
		} else {
			keyObj := memory.Assets.Rooms[memory.User.Room].GetObjectKey(argA)
			if keyObj == "" {
				printWarning(fmt.Sprintf("There's no object named %s here.\n", argA))
			} else {
				obj := memory.Assets.GetObject(keyObj)
				fmt.Printf("%s\n", obj.Description)
			}
		}
	case "inventory":
		memory.User.ListInventory()
	case "help":
		fmt.Println("\x1b[35mPossible commands:")
		fmt.Println("  use OBJECT -> interact with a scene object")
		fmt.Println("  use OBJECT with ITEM -> interact with OBJECT using ITEM from inventory")
		fmt.Println("  get OBJECT -> puts OBJECT in your inventory")
		fmt.Println("  check OBJECT -> get OBJECT description")
		fmt.Println("  inventory -> list what is in your inventory")
		fmt.Println("  help -> command list")
		fmt.Println("  save FILENAME -> save current game to FILENAME")
		fmt.Println("  load FILENAME -> load game stored in FILENAME")
		fmt.Println("  newgame -> start a new game\x1b[0m")
	case "save":
		if argA == "" {
			printError("Invalid command.")
		} else {
			memory.SaveData(argA)
		}
	case "load":
		if argA == "" {
			printError("Invalid command.")
		} else {
			if ans, usr := memory.LoadData(argA); ans {
				memory.User = usr
				r = memory.User.Room
				memory.User.Room = -1
				printNextWait()
			}
		}
	case "newgame":
		fmt.Print("\x1b[31mAre you sure? (Y/N)\x1b[0m ")
		var ans string
		fmt.Scanln(&ans)
		if ans == "Y" || ans == "y" {
			r = -2
		}
	case "quit":
		fmt.Print("\x1b[31mAre you sure? (Y/N)\x1b[0m ")
		var ans string
		fmt.Scan(&ans)
		if ans == "Y" || ans == "y" {
			r = -1
		}
	case "":
	default:
		printError(fmt.Sprintf("There is no command %s. Type 'help' for command list.", cmd))
	}

	return r
}

func Start() {
	fmt.Print("\x1b[0m")
	memory := createMemory()
	nroom, proom := 0, -1
	for {
		nroom = loop(proom != memory.User.Room, memory)
		if nroom == -1 {
			break
		} else if nroom == -2 {
			memory.User.Room = -1
			nroom = 0
			memory.User.Inventory = make(map[string]string)
			memory.User.Solved = make([]string, 0)
		}

		proom = memory.User.Room
		memory.User.Room = nroom
	}
}
