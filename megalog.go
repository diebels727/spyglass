package megalog

import(
  "github.com/diebels727/tangent"
)

type MegaLog struct {
  Nick,Username,Server,Channel string
  Port int
  Context *tangent.Context
}

func New(nick string,username string,server string,port int,channel string) (*MegaLog) {
  context := tangent.New()
  return &MegaLog{nick,username,server,channel,port,context}
}




func (m *MegaLog) Run() {


  //m.Connect()
  // m.Identify()
  // go func() {
  //   message <- log
  //   //write log
  // }
  // m.JoinChannel()
  // go func() {
  //   m.Beacon()
  // }
}