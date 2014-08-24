package spyglass

import ("net"
        "log"
        "bufio"
        "fmt"
        "net/textproto"
        "strings"
      )

type Bot struct{
        server string
        port string
        nick string
        user string
        pass string
        display,write,ping chan string
        events chan *Event
        conn net.Conn
        eventHandlers map[string]func(event *Event)
}

type Event struct {
  Source string
  Command string
  RawCommand string
  Arguments string
  RawMessage string
  RawArguments string
}

//will become spyglass/event
func EventNew(message string) (*Event) {
  e := &Event{RawMessage: message}
  e.Parse()
  return e
}

func (e *Event) Parse() {
  message := e.RawMessage
  current_message := e.RawMessage

  if message[0:1] == ":" {
    if i := strings.Index(message," "); i > -1 {
      current_message = message[i+1:len(message)] //peel off source
      e.Source = message[0:i]
    } else {
      log.Println("Server IRC protocol error.  Expected :<source> CMD ARGS, got ",message)
    }
  }

  message = current_message
  current_message = message
  if i := strings.Index(message," "); i > -1 {
    current_message = message[i+1:len(current_message)]
    e.RawCommand = message[0:i]
    e.Command = strings.ToUpper(e.RawCommand)
    e.RawArguments = message[i+1:]
  } else {
    log.Println("Server IRC protocol error. Expected CMD ARGS, got",message)
  }

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

    //decompose messages
    //messages are in this format:
    //
    // :<source> COMMAND <ARGS>
    //
    // :<source> is optional

    event := EventNew(line)
    bot.events <- event
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

// type handler func()

func (bot *Bot) RegisterEventHandler(command string,handler func(event *Event) ) {
  bot.eventHandlers[command] = handler
}

func (bot *Bot) handleEvent(event *Event) {
  event_handler := bot.eventHandlers[event.Command]
  if event_handler != nil {
    event_handler(event)
  }
}

func (bot *Bot) Run() {
  if bot.conn != nil {
    //defend against running a dupe
  }

  display := make(chan string,1024)
  write := make(chan string,1024)
  ping := make(chan string,1)
  events := make(chan *Event,1024)

  bot.display = display
  bot.write = write
  bot.ping = ping
  bot.events = events


  bot.eventHandlers = make(map[string]func(event *Event))

  bot.RegisterEventHandler("PING",func(event *Event) {
    fmt.Println("[Event PING]: Handling ",event)
    bot.write <- fmt.Sprintf("PONG %s\r\n",event.Source)
  })


  //display loop
  go func() {
    for {
      message := <- bot.display
      fmt.Println(message)  // switch this around to a io.Writer obj; probably logger interface
    }
  }()

  //event loop
  go func() {
    for {
      event := <- bot.events
      bot.handleEvent(event)
    }
  }()

  //write loop
  go func() {
    bot.WriteLoop()
  }()

  // go func() {
  //   for {
  //     select {
  //       case message := <- bot.ping:
  //         fmt.Println("Ping Received, responding...")
  //         bot.write <- fmt.Sprintf("PONG %s\r\n",message)
  //         fmt.Println("Responded with pong...")
  //       default:
  //         //if timer, then ping
  //     }
  //   }
  // }()

  //read loop
  go func() {
    reader := bufio.NewReader(bot.conn)
    tp := textproto.NewReader(reader)
    bot.ReadLoop(tp)
  }()

  ready <- true
}
