package game

import (
	"encoding/json"
	"fmt"
	"os"
)

type Assets struct {
	Objects []Object
	Rooms   []Room
}

func (d Assets) GetObject(key string) *Object {
	for _, o := range d.Objects {
		if o.Key == key {
			return &o
		}
	}
	return nil
}

type Memory struct {
	Assets Assets
	User   User
}

type User struct {
	Inventory map[string]string
	Solved    []string
	Room      int
}

func (u User) ListInventory() {
	if len(u.Inventory) == 0 {
		fmt.Println("\x1b[94mNothing to see here.\x1b[0m")
	} else {
		fmt.Print("\x1b[94m")
		for k := range u.Inventory {
			fmt.Printf("%s\n", k)
		}
		fmt.Print("\x1b[0m")
	}
}

// func (u User) AddSolved(object string) {
// 	if !u.IsSolved(object) {
// 		u.Solved = append(u.Solved, object)
// 	}
// }

func (u User) IsSolved(object string) bool {
	if object == "" {
		return true
	}
	for _, v := range u.Solved {
		if v == object {
			return true
		}
	}
	return false
}

func (u User) GetInventoryItem(item string) string {
	for k, v := range u.Inventory {
		if k == item {
			return v
		}
	}
	return ""
}

func LoadAssets() Assets {
	fdata, _ := os.ReadFile("data.json")
	var data Assets
	json.Unmarshal(fdata, &data)
	return data
}

func (m Memory) SaveData(file string) {
	f, e := os.Stat(file)
	if os.IsExist(e) {
		if f.IsDir() {
			fmt.Println("\x1b[1;31m'" + file + "' is a directory.\x1b[0m")
		} else {
			fmt.Println("\x1b[1;31mThere's already a file named '" + file + "'. Overwite? (Y/N).\x1b[0m")
			var ans string
			fmt.Scan(&ans)
			if ans != "Y" && ans != "y" {
				return
			}
		}
	}

	save, _ := json.Marshal(m.User)
	if os.IsExist(e) {
		os.Remove(file)
	}
	os.WriteFile(file, save, 0644)
	fmt.Println("\x1b[1;32mData saved successfully to '" + file + "'.\x1b[0m")
}

func (m Memory) LoadData(file string) (bool, User) {
	var usr User
	f, e := os.Stat(file)
	if os.IsNotExist(e) || f.IsDir() {
		fmt.Println("\x1b[1;31mThe file '" + file + "' does not exist.\x1b[0m")
	} else {
		save, _ := os.ReadFile(file)
		if !json.Valid(save) {
			fmt.Println("\x1b[1;31mThe file '" + file + "' insn't a valid save data.\x1b[0m")
		} else {
			printWarning("Current progress will be lost if not saved. Continue? (Y/N)")
			var ans string
			fmt.Scanln(&ans)
			if ans == "Y" || ans == "y" {
				json.Unmarshal(save, &usr)
				fmt.Println("\x1b[1;32mData loaded successfully.\x1b[0m")
				return true, usr
			}
		}
	}
	return false, usr
}
