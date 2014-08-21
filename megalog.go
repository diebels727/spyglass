package megalog

import(
  "github.com/diebels727/tangent"
  irc "github.com/xales/tangent/irc"
)

type MegaLog struct {
  User *tangent.User
  Server,Channel,Port string
  Context *tangent.Context
}

func New(nick string,username string,server string,port string,channel string) (*MegaLog) {
  context := tangent.New()
  hostmask := irc.HostMask{Nick: nick,Ident: username}
  user := &tangent.User{HostMask: hostmask, Real: username}
  return &MegaLog{ User:user,
                   Server: server,
                   Channel: channel,
                   Port: port,
                   Context: context }
}

func (m *MegaLog) Connect() {
  context := m.Context
  context.Connect(m.Server+":"+m.Port,false,m.User)

}

func (m *MegaLog) Run() {
  m.Connect()
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