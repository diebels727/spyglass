package spyglass

import(
  "testing"
  "time"
  "fmt"
)

//just a basic sniff test
func TestInit(t *testing.T) {
  ready = make(chan bool,1)
  m := New("localhost","6667","nick89122","jiggly101001","")
  fmt.Println("[Run] Connecting")
  conn := m.Connect()
  defer conn.Close()

  fmt.Println("[Run] Connect finished execution.")
  m.Run()


  user_cmd := fmt.Sprintf("USER %s 8 * :%s\r\n", m.nick, m.nick)
  nick_cmd := fmt.Sprintf("NICK %s\r\n", m.nick)
  fmt.Println("[TestInit] Sending USER command")
  m.write <- user_cmd
  fmt.Println("[TestInit] Sending NICK command")
  m.write <- nick_cmd
  fmt.Println("[TestInit] Sending JOIN command")
  m.Join("#cinch-bots")
  fmt.Println("[TestInit] Sending JOIN command")
  m.write <- "JOIN #foofoo\r\n"

  fmt.Println("Sleeping for 5 seconds")
  time.Sleep(time.Second * 5)

  for {
    time.Sleep(time.Second * 1)
  }
}
