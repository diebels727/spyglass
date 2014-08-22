package spyglass 

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
        pass string
        display,write,ping chan string
        conn net.Conn
}

var ready chan bool

func New(server string,port string,nick string,user string,pass string) *Bot {
  return &Bot{server: server,
              port: port,
              nick: nick,
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

    event := line
    //if line is an event, dispatch to handle
    fmt.Println("LINE[0]: ",event[0:1]==":")
    if line[0:4] == "PING" {
      fmt.Println("PING Received!")
      message := event[5:]
      fmt.Println("[PING]:",message)
      bot.ping <- message
    }
    //bot.read is bot.display, change this

    //strip /r/n from received messages
    //convert line to an event
    //decode the event
    //if begins with :, then handle with msg type
    //if begins with PING, then handle with ping loop

    bot.display <- line
  }
}

func (bot *Bot) WriteLoop() {
  for {
    select {
      case cmd := <- bot.write:
        fmt.Println("[WriteLoop] Command: ",cmd)
        bot.RawCmd(cmd)
      default:
    }
  }
}

func (bot *Bot) Run() {
  if bot.conn != nil {
    //defend against running a dupe
  }

  display := make(chan string,1024)
  write := make(chan string,1024)
  ping := make(chan string,1)

  bot.display = display
  bot.write = write
  bot.ping = ping

  //display loop
  go func() {
    for {
      message := <- bot.display
      fmt.Println(message)  // switch this around to a io.Writer obj; probably logger interface
    }
  }()

  //write loop
  go func() {
    bot.WriteLoop()
  }()

  go func() {
    for {
      select {
        case message := <- bot.ping:
          fmt.Println("Ping Received, responding...")
          bot.write <- fmt.Sprintf("PONG %s\r\n",message)
          fmt.Println("Responded with pong...")
        default:
          //if timer, then ping
      }
    }
  }()

  //read loop
  go func() {
    reader := bufio.NewReader(bot.conn)
    tp := textproto.NewReader(reader)
    bot.ReadLoop(tp)
  }()

  ready <- true
}
