package megalog

import(
  "testing"
  "time"
  "fmt"
)

//just a basic sniff test
func TestInit(t *testing.T) {
  ready = make(chan bool,1)
  m := New("irc.freenode.net","6667","nick89122","jiggly101001","#cinch-bots","")
  m.Run()

  <- ready

  user_cmd := fmt.Sprintf("USER %s 8 * :%s\r\n", m.nick, m.nick)
  nick_cmd := fmt.Sprintf("NICK %s\r\n", m.nick)
  join_cmd := fmt.Sprintf("JOIN %s\r\n", m.channel)
  m.write <- user_cmd
  m.write <- nick_cmd
  m.write <- join_cmd
  m.write <- "JOIN #foofoo\r\n"



  for {
    time.Sleep(time.Second * 1)
  }
}