package daemon

import (
	"fmt"
	"os"
	"os/exec"
)

func Run() {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "-d=true" {
			args[i] = "-d=false"
			break
		}
	}

	cmd := exec.Command(os.Args[0], args...)
	cmd.Start()
	fmt.Println("[PID]", cmd.Process.Pid)
	os.Exit(0)
}
