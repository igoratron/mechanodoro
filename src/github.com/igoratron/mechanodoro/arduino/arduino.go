package arduino

import "log"
import "time"
import "github.com/tarm/serial"

var FLAG_RAISED = byte(255)
var FLAG_LOWERED = byte(128)

type Arduino struct {
  Name string
  connection *serial.Port
}

func (a *Arduino) Start() {
  c := &serial.Config{Name: a.Name, Baud: 115200}
  var err error
  a.connection, err = serial.OpenPort(c)
  if err != nil {
    log.Fatal("Can't talk to the Arduino. ", err)
  }

  waitForArduino(a.connection)
}

func (a *Arduino) Stop() {
  a.connection.Close()
}

func (a *Arduino) send(b byte) {
  _, err := a.connection.Write([]byte{b})
  if err != nil {
    log.Fatal("Failed to send to Arduino", err)
  }
}

func (a *Arduino) RaiseFlag() {
  a.send(FLAG_RAISED)
}

func (a *Arduino) LowerFlag() {
  a.send(FLAG_LOWERED)
}

func waitForArduino(port *serial.Port) {
  var isConnectionEstablished bool

  log.Printf("Establishing connection with the Arduino")
  defer log.Println("Connected to Arduino")

  go func() {
    for !isConnectionEstablished {
      _, err := port.Write([]byte{0})
      if err != nil {
        log.Println(err)
      }
      time.Sleep(time.Until(time.Now().Add(200 * time.Millisecond)))
    }
  }()

  buf := make([]byte, 32)
  _, err := port.Read(buf)
  if err != nil {
    log.Fatal(err)
  }
  isConnectionEstablished = true
}
