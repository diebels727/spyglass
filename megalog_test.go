package megalog

import "testing"

//just a basic sniff test
func TestInit(t *testing.T) {
  m := New("irc.freenode.net","6667","nick89122","jiggly101001","#cinch-bots","")
  m.Run()
  m.Join("#foofoo")
}