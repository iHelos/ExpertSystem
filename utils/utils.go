package utils

import (
	"os"
	"os/exec"
	"os/signal"
)

func CheckInterrupt(f func ()){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(){
		for range c {
			f()
			os.Exit(0)
		}
	}()
}

func CLS(){
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}