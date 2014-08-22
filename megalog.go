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
  fmt.Println("[RawCmd] Raw command message: ",message)
  fmt.Fprintf(bot.conn,message)
}

func (bot *Bot) ReadLoop(tp *textproto.Reader) {
  for {
    fmt.Println("[ReadLoop] Reading from textproto.")
    line, err := tp.ReadLine()
    fmt.Println("[ReadLoop] Read from textproto.")

    if err != nil {
      fmt.Println("[ReadLoop] Error: ",err)
      break // break loop on errors
    }
    //if line is an event, dispatch to handle
    fmt.Println("[ReadLoop] Raw line.")
    bot.read <- line
  }
}

func (bot *Bot) WriteLoop() {
  for {
    select {
      case cmd := <- bot.write:
        fmt.Println("[WriteLoop] Received command.")
        bot.RawCmd(cmd)
        fmt.Println("[WriteLoop] Sent command.")
      default:
    }
  }
}

func (bot *Bot) Run() {
  fmt.Println("[Run] Connecting")
  // conn := bot.Connect()
  bot.Connect()

  fmt.Println("[Run] Connect finished execution.")

  // defer conn.Close()

  read := make(chan string,1024)
  write := make(chan string,1024)

  bot.read = read
  bot.write = write

  //display loop
  go func() {
    fmt.Println("[Run] Launching DisplayLoop")
    for {
      fmt.Println("[Run] Preparing to read from read channel.")
      message := <- bot.read
      fmt.Println("[Run] Read from channel.")

      fmt.Println(message)  // switch this around to a io.Writer obj; probably logger interface
    }
  }()

  //write loop
  go func() {
    fmt.Println("[Run] Launching WriteLoop")
    bot.WriteLoop()
  }()

  //read loop
  go func() {
    fmt.Println("[Run] Launching ReadLoop")
    reader := bufio.NewReader(bot.conn)
    fmt.Println("[Run] Checking bot.conn: ",bot.conn)
    tp := textproto.NewReader(reader)
    bot.ReadLoop(tp)
  }()

  fmt.Println("[Run] Triggering ready.")
  ready <- true
  fmt.Println("[Run] Triggered ready.")

}
