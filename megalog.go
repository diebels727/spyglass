package megalog

import ("net"
        "log"
        "bufio"
        "fmt"
        "net/textproto"
        "time"
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
  fmt.Fprintf(bot.conn,"JOIN %s\r\n",channel)
}

func (bot *Bot) RawCmd(message string) {
  fmt.Fprintf(bot.conn,message)
}

func (bot *Bot) Loop() {
  reader := bufio.NewReader(bot.conn)
  tp := textproto.NewReader( reader )
  for {
    line, err := tp.ReadLine()
    if err != nil {
      break // break loop on errors
    }
    fmt.Println(line)
    // bot.read <- line
    select {
    case command := <- bot.write:
      bot.RawCmd(command)
      fmt.Println("sent message: ", command)
    default:
      fmt.Println("none")
    }
  }
}

func (bot *Bot) Run() {
  conn := bot.Connect()
  defer conn.Close()

  read := make(chan string)
  write := make(chan string)

  bot.read = read
  bot.write = write

  go func(read chan string) {
    message := <- bot.read
    fmt.Println(message)
  }(read)

  go func() {
    bot.Loop()
  }()

  user_cmd := fmt.Sprintf("USER %s 8 * :%s\r\n", bot.nick, bot.nick)
  nick_cmd := fmt.Sprintf("NICK %s\r\n", bot.nick)
  join_cmd := fmt.Sprintf("JOIN %s\r\n", bot.channel)

  bot.write <- user_cmd
  bot.write <- nick_cmd
  bot.write <- join_cmd

  for {
    time.Sleep(time.Second * 1)
  }

}
