package megalog

import "testing"

//just a basic sniff test
func TestInit(t *testing.T) {
  m := New("nick","username","server","6667","channel")
  if m.Server != "server" {
    t.Error("Expected \"server\" got ",m.Server)
  }
}

func TestRun(t *testing.T) {
  m := New("merryNicksMas","sudormrf","irc.freenode.org","6667","#cinch-bots")
  m.Run()
}