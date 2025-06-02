package main

import (
	"fmt"
	"os/exec"
)

func runBastille(args ...string) error {
	cmd := exec.Command("bastille", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bastille %v failed: %v\n%s", args, err, out)
	}
	return nil
}

