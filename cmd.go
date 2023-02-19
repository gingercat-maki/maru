package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
		cmd, err := exec.Command("echo", "echo maki").Output()
		if err != nil {
			fmt.Printf("[%v] error is: %s \n", i, err)
		} else {
			output := string(cmd)
			fmt.Printf("[%v] error is: %s \n", i, output)
		}
	}
}
