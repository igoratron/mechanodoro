package pomodoro

import (
 "time"
)


type Pomodoro struct {
  pomodoroStart time.Time
  isRunning bool
}

func (p *Pomodoro) Start() time.Time {
  p.pomodoroStart = time.Now()
  p.isRunning = true

  return p.pomodoroStart
}

func (p *Pomodoro) GetCurrentDuration() time.Duration {
  if !p.isRunning { return time.Duration(0) }
  return time.Since(p.pomodoroStart)
}

func (p *Pomodoro) Stop() {
  p.isRunning = false
}
