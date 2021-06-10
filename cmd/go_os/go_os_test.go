package go_os

import (
	"fmt"
	"os"
	"testing"
)

func TestHostname(t *testing.T) {
	hostname, _ := os.Hostname()
	fmt.Printf("%s\n", hostname)
}
