package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	shellquote "github.com/kballard/go-shellquote"
)

func main() {
	commands := os.Args[1:]
	wg := sync.WaitGroup{}
	for _, command := range commands {
		var words []string
		if command == "caddy" {
			words = []string{"caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"}
		} else {
			var err error
			words, err = shellquote.Split(command)
			if err != nil {
				log.Printf("invalid command: %s: %v", command, err)
				continue
			}
		}

		wg.Add(1)
		go func() {
			if err := run(words); err != nil {
				log.Printf("child error: %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func run(arg []string) error {
	cmd := exec.Command(arg[0], arg[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	wait := make(chan error, 1)
	go func() {
		wait <- cmd.Wait()
		close(wait)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig)

	for {
		select {
		case sig := <-sig:
			if err := cmd.Process.Signal(sig); err != nil {
				log.Print("child signal error", sig, err)
			}
		case err := <-wait:
			var waitStatus syscall.WaitStatus
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus = exitError.Sys().(syscall.WaitStatus)
				os.Exit(waitStatus.ExitStatus())
			}
			if err != nil {
				return err
			}
			return nil
		}
	}
}
