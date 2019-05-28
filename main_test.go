package main

import (
	"os/exec"
	"regexp"
	"testing"
)

func TestHelpOptionReturnsSuccessAndPrintsToStdOut(t *testing.T) {
	out, err := exec.Command("go", "run", "./main.go", "-help").Output()
	if err != nil {
		t.Fatalf("go run error: %#v", err)
	}
	re := regexp.MustCompile(`Usage of`)
	if re.Match(out) == false {
		t.Fatalf("output does not include help text: %s", string(out))
	}
}
