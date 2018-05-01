package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// returns available memory in kilobytes
func avail() (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("bash", "-c", "dmesg | grep memory")
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.Split(strings.Split(buf.String(), "\n")[1], " ")[3][:8], nil
}

// returns free memory in kilobytes
func free() (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("vmstat", "-H")
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.Split(strings.Split(buf.String(), "\n")[2], " ")[4], nil
}

type Unit int8

const (
	KB Unit = iota
	MB
	GB
)

// convert string representation to float
func conv(s string, u Unit) float32 {
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}

	if u == KB {
		return float32(f64)
	} else if u == MB {
		return float32(f64) / 1000
	} else if u == GB {
		return float32(f64) / 1000 / 1000
	}

	return float32(f64)
}

func main() {
	avail, err := avail()
	if err != nil {
		panic(err)
	}

	free, err := free()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%0.2f available\n", conv(avail, GB))
	fmt.Printf("%0.2f free\n\n", conv(free, GB))

	fmt.Printf("%0.2f%% used\n",
		100*(conv(avail, KB)-conv(free, KB))/conv(avail, KB))
}
