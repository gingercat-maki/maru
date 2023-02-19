package main

import (
	"fmt"
	"os/exec"
	"time"
)

// tctl --namespace maki-test wf start --tq temporal-bench --wt bench-workflow --wtt 5 --et 1800 --if ./scenarios/basic-test.json --w

func main() {

	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
		cmd, err := exec.Command("echo", "maki").Output()
		if err != nil {
			fmt.Println("[%d] error is: %s", i, err)
		} else {
			output := string(cmd)
			fmt.Println("[%d] output is: %s", i, output)
		}
	}

}
