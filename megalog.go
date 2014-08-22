package megalog

import ("net"
        "log"
        "bufio"
        "fmt"
        "net/textproto"
      )

type Bot struct{
        server string
        port string
        nick string
        user string
        channel string
        pass string
        read,write chan string
        conn net.Conn
}

var ready chan bool

func New(server string,port string,nick string,user string,channel string,pass string) *Bot {
  return &Bot{server: server,
              port: port,
              nick: nick,
              channel: channel,
              pass: "",
              conn: nil,
              user: user}
}

func (bot *Bot) Connect() (conn net.Conn){
  connection_string := fmt.Sprintf("%s:%s",bot.server,bot.port)
  conn, err := net.Dial("tcp",connection_string)
  if err != nil{
    log.Fatal("unable to connect to IRC server ", err)
  }
  bot.conn = conn
  log.Printf("Connected to IRC server %s (%s)\n", bot.server, bot.conn.RemoteAddr())
  return bot.conn
}

func (bot *Bot) Join(channel string) {
  bot.write <- fmt.Sprintf("JOIN %s\r\n",channel)
}

func (bot *Bot) RawCmd(message string) {
  fmt.Fprintf(bot.conn,message)
}

func (bot *Bot) ReadLoop(tp *textproto.Reader) {
  for {
    line, err := tp.ReadLine()
    if err != nil {
      break // break loop on errors
    }
    //if line is an event, dispatch to handle
    bot.read <- line
  }
}

func (bot *Bot) WriteLoop() {
  for {
    select {
      case cmd := <- bot.write:
        bot.RawCmd(cmd)
      default:
    }
  }
}

func (bot *Bot) Run() {
  conn := bot.Connect()

  defer conn.Close()

  read := make(chan string,1024)
  write := make(chan string,1024)

  bot.read = read
  bot.write = write

  //display loop
  go func() {
    for {
      message := <- bot.read
      log.Println(message)  // switch this around to a io.Writer obj; probably logger interface
    }
  }()

  //write loop
  go func() {
    bot.WriteLoop()
  }()

  //read loop
  go func() {
    reader := bufio.NewReader(bot.conn)
    tp := textproto.NewReader( reader )
    bot.ReadLoop(tp)
  }()

  ready <- true

}
