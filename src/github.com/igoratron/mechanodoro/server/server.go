package server

import (
 "log"
 "net"
 "sync"
 "time"
 "bufio"
)

const SOCKET = "/tmp/mechanodoro.sock"

type Server struct {
  isRunning bool
  openConnections sync.WaitGroup
  socket net.Listener
  Commands map[string]func() string
  OnClientConnected func(func(string) error, *bool)
}

func (s *Server) Start() {
  var err error
  s.socket, err = net.Listen("unix", SOCKET)

  if err != nil {
    log.Fatal(err)
  }

  s.isRunning = true
  log.Println("Running")

  go handleConnections(s)
}

func (s *Server) Stop() {
  s.isRunning = false
  s.openConnections.Wait()
  s.socket.Close()
}

func handleConnections(s *Server) {
  for s.isRunning {
    connection, err := s.socket.Accept()

    if err != nil && s.isRunning {
      log.Println("Accept error", err)
    } else {
      s.openConnections.Add(1)
      go clientConnected(s, &connection)
    }
  }
}

func clientConnected(s *Server, connection *net.Conn) {
  log.Println("Client connected")

  defer (*connection).Close()
  defer s.openConnections.Done()

  for s.isRunning {
    scanner := bufio.NewScanner(*connection)
    (*connection).SetDeadline(time.Now().Add(time.Second))

    for scanner.Scan() {
      command := scanner.Text()
      getResponse, ok := s.Commands[command]

      if ok {
        (*connection).Write([]byte(getResponse() + "\n"))
      } else {
        (*connection).Write([]byte("Unrecognised command"))
      }
    }

    if err, ok := scanner.Err().(net.Error); !ok || !err.Timeout() {
      break
    }
  }

  log.Println("Client disconnected")
}
