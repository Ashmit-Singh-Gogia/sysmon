package main

import (
	"fmt"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/proc"
)

func main() {
	snap, err := proc.ReadStat()
	if err != nil {
		fmt.Println("There is some error in main file")
	}
	fmt.Println(snap.Total.ID)
}
