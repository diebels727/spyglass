package spyglass

import ("net"
        "log"
        "bufio"
        "fmt"
        "net/textproto"
        "strings"
        "time"
      )

type Bot struct{
  server string
  port string
  nick string
  user string
  pass string
  display,write,ping chan string
  events chan *Event
  datastore chan *Event
  Ready,Stopped chan bool
  Conn net.Conn
  eventHandlers map[string]func(event *Event)

  Datastore Datastore
}

type Datastore interface {
  Write(event *Event)
}

type Event struct {
  Source string
  User string
  URI string
  Command string
  RawCommand string
  Arguments string
  RawMessage string
  RawArguments string
  Target string
  Message string
  Timestamp int64
}

//will become spyglass/event
func EventNew(message string) (*Event) {
  e := &Event{RawMessage: message}
  e.Parse()
  return e
}

func (e *Event) Parse() {
  e.Timestamp = time.Now().Unix()

  message := e.RawMessage
  current_message := e.RawMessage

  //guard against an anomaly where the message length is zero
  if (len(message) >= 1) && (message[0:1] == ":") {
    if i := strings.Index(message," "); i > -1 {
      current_message = message[i+1:len(message)] //peel off source
      e.Source = message[1:i]
      source := e.Source
      if (len(source) >= 1) {
        if i := strings.Index(source,"~!"); i > -1 {
          e.User = source[0:i]
          e.URI = source[i+1:len(source)]
        }
      }

    } else {
      log.Println("Server IRC protocol error.  Expected :<source> CMD ARGS, got ",message)
    }
  }



  message = current_message
  // current_message = message

  if i := strings.Index(message," "); i > -1 {
    current_message = message[i+1:len(current_message)]
    e.RawCommand = message[0:i]
    e.Command = strings.ToUpper(e.RawCommand)
    e.RawArguments = message[i+1:]

    if j := strings.Index(e.RawArguments," ");j > -1 {
      e.Target = e.RawArguments[0:j]

      if k := strings.Index(e.RawArguments[j+1:],":");k > -1 {
        e.Message = e.RawArguments[j+k+1:]
      }
    }

  } else {
    log.Println("Server IRC protocol error. Expected CMD ARGS, got",message)
  }
}


func New(server string,port string,nick string,user string,pass string) *Bot {
  bot := &Bot{server: server,
              port: port,
              nick: nick,
              pass: "",
              Conn: nil,
              user: user}
  bot.eventHandlers = make(map[string]func(event *Event))
  return bot
}

func (bot *Bot) Connect() (conn net.Conn){
  connection_string := fmt.Sprintf("%s:%s",bot.server,bot.port)
  conn, err := net.Dial("tcp",connection_string)
  if err != nil{
    log.Fatal("unable to connect to IRC server ", err)
  }
  bot.Conn = conn
  log.Printf("Connected to IRC server %s (%s)\n", bot.server, bot.Conn.RemoteAddr())
  return bot.Conn
}

func (bot *Bot) GetNick() string {
  return bot.nick
}

func (bot *Bot) Join(channel string) {
  bot.write <- fmt.Sprintf("JOIN %s\r\n",channel)
}

func (bot *Bot) User() {
  bot.write <- fmt.Sprintf("USER %s 8 * :%s\r\n",bot.nick,bot.user)
}

func (bot *Bot) Nick() {
  bot.write <- fmt.Sprintf("NICK %s\r\n",bot.nick)
}

func (bot *Bot) List() {
  bot.write <- fmt.Sprintf("LIST\r\n")
}

//directly write to the server's connection, bypassing all scheduling.
func (bot *Bot) RawCmd(message string) {
  fmt.Fprintf(bot.Conn,message)
}

func (bot *Bot) Cmd(message string) {
  bot.write <- fmt.Sprintf("%s\r\n",message)
}

func (bot *Bot) Send(message string) {
  bot.write <- message
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
    bot.datastore <- event
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
  bot.Ready = make(chan bool,1)
  bot.Stopped = make(chan bool,1)

  if bot.Conn != nil {
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

  bot.RegisterEventHandler("PING",func(event *Event) {
    bot.write <- fmt.Sprintf("PONG %s\r\n",event.RawArguments)
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
      bot.handleEvent(event)  //perhaps a heap, and assign priority to events?
    }
  }()

  //logging to datastore support
  datastore := make(chan *Event,1024)
  bot.datastore = datastore
  // datastore loop
  go func() {
    for {
      event := <- bot.datastore
      bot.Datastore.Write(event)
    }
  }()

  //write loop
  go func() {
    bot.WriteLoop()
  }()

  //read loop
  go func() {
    reader := bufio.NewReader(bot.Conn)
    tp := textproto.NewReader(reader)
    bot.ReadLoop(tp)
  }()

  bot.Ready <- true
}
