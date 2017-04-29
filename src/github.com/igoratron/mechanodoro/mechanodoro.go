package main

import (
 "log"
 "os"
 "os/signal"
 "syscall"
 "time"
 "fmt"
)

import "github.com/igoratron/mechanodoro/server"
import "github.com/igoratron/mechanodoro/pomodoro"


var pomo pomodoro.Pomodoro

func main() {
  pomo = pomodoro.Pomodoro{}
  server := server.Server{
    Commands: map[string]func() string {
      "start-task": func() string {
        pomo.Start()
        return "OK: Task started"
      },
      "start-short-break": func() string {
        pomo.Start()
        return "OK: Short break on"
      },
      "start-long-break": func() string {
        pomo.Start()
        return "OK: Long break on"
      },
      "stop": func() string {
        pomo.Stop()
        return "OK: Pomodoro stopped"
      },
      "get-duration": func() string {
        return humanizeDuration(pomo.GetCurrentDuration())
      },
    },
  }
  server.Start()


  s := <- waitForInterrupt()
  log.Printf("Got signal %s, exitting...\n", s)
  server.Stop()
}

func waitForInterrupt() chan os.Signal {
  interruptChannel := make(chan os.Signal, 1)
  signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
  return interruptChannel
}

func humanizeDuration(d time.Duration) string {
  return fmt.Sprintf("%02d:%02d", int(d.Minutes()), int(d.Seconds()) % 60)
}
