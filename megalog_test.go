package megalog

import(
  "testing"
  "time"
  "fmt"
)

//just a basic sniff test
func TestInit(t *testing.T) {
  ready = make(chan bool,1)
  m := New("localhost","6667","nick89122","jiggly101001","#cinch-bots","")
  m.Run()

  // <- ready
  fmt.Println("[TestInit] Sleeping for 15 seconds")
  time.Sleep(time.Second * 15)
  fmt.Println("[TestInit] Done sleeping.")


  user_cmd := fmt.Sprintf("USER %s 8 * :%s\r\n", m.nick, m.nick)
  nick_cmd := fmt.Sprintf("NICK %s\r\n", m.nick)
  join_cmd := fmt.Sprintf("JOIN %s\r\n", m.channel)
  fmt.Println("[TestInit] Sending USER command")
  m.write <- user_cmd
  fmt.Println("[TestInit] Sending NICK command")
  m.write <- nick_cmd
  fmt.Println("[TestInit] Sending JOIN command")
  m.write <- join_cmd
  fmt.Println("[TestInit] Sending JOIN command")
  m.write <- "JOIN #foofoo\r\n"

  for {
    time.Sleep(time.Second * 1)
  }
}