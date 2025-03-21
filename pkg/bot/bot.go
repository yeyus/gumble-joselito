package bot

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yeyus/gumble-joselito/pkg/joselito"
	"github.com/yeyus/gumble/gumble"
	"github.com/yeyus/gumble/gumbleutil"
	_ "github.com/yeyus/gumble/opus"
)

type StreamerState int

const (
	StreamerStateDisconnected StreamerState = iota
	StreamerStateConnecting
	StreamerStateConnected
	StreamerStateIdle
	StreamerStateTalking
)

type Bot struct {
	logger *log.Logger
	State  StreamerState

	// Mumble section
	Config *gumble.Config
	Client *gumble.Client

	MumbleAddress string
	TLSConfig     *tls.Config

	Room []string

	// Joselito section
	wsEndpoint          string
	wsHeaders           http.Header
	websocketConnection *websocket.Conn
	session             *joselito.Session
	stream              *joselito.Stream
	talkgroups          []*joselito.DMRID

	WaitGroup *sync.WaitGroup
}

func NewBot(address string, room []string, config *gumble.Config, tlsConfig *tls.Config, wsEndpoint string, wsHeaders http.Header, talkgroups []*joselito.DMRID) *Bot {
	return &Bot{
		logger:        log.New(os.Stdout, "[bot] ", log.LstdFlags),
		State:         StreamerStateDisconnected,
		Config:        config,
		TLSConfig:     tlsConfig,
		MumbleAddress: address,
		Room:          room,
		WaitGroup:     new(sync.WaitGroup),
		wsEndpoint:    wsEndpoint,
		wsHeaders:     wsHeaders,
		talkgroups:    talkgroups,
	}
}

func (s *Bot) Connect() error {
	if s.State != StreamerStateDisconnected {
		return nil
	}

	s.State = StreamerStateConnecting
	s.WaitGroup.Add(1)

	s.Config.Attach(gumbleutil.Listener{
		Connect:    s.onConnect,
		Disconnect: s.onDisconnect,
		UserChange: s.onUserChange,
	})

	client, err := gumble.DialWithDialer(new(net.Dialer), s.MumbleAddress, s.Config, s.TLSConfig)
	if err != nil {
		s.State = StreamerStateDisconnected
		s.logger.Printf("connect: error while connecting %v\n", err)
		return err
	}

	s.Client = client
	return nil
}

func (s *Bot) onConnect(e *gumble.ConnectEvent) {
	s.State = StreamerStateConnected
	s.logger.Printf("connected to %s", s.MumbleAddress)

	targetChannel := e.Client.Channels.Find(s.Room...)
	if targetChannel == nil {
		s.logger.Printf("could not find channel %s, aborting", s.Room)
		e.Client.Disconnect()
		return
	}

	s.logger.Printf("moving to %s\n", targetChannel.Name)
	e.Client.Self.Move(targetChannel)
}

func (s *Bot) Disconnect() {
	s.Client.Disconnect()
}

func (s *Bot) onDisconnect(e *gumble.DisconnectEvent) {
	defer s.WaitGroup.Done()

	s.State = StreamerStateDisconnected
	s.logger.Printf("onDisconnect: disconnected from %s", s.MumbleAddress)
}

func (s *Bot) onUserChange(e *gumble.UserChangeEvent) {
	if e.Type == gumble.UserChangeChannel {
		// users connected to our channel
		numUsersInChannel := len(e.Client.Self.Channel.Users)
		s.logger.Printf("onUserChange: users in our channel %d!", numUsersInChannel)

		if numUsersInChannel > 1 && s.State == StreamerStateConnected {
			// someone arrived
			s.State = StreamerStateIdle
			s.logger.Printf("onUserChange: someone is in the channel, new state is %v\n", s.State)
			s.StartStreaming()
		} else if numUsersInChannel <= 1 && s.State == StreamerStateIdle {
			// I'm alone here
			s.State = StreamerStateConnected
			s.logger.Printf("onUserChange: everyone left the room, new state is %v\n", s.State)
			s.StopStreaming()
		}
	}
}

func (s *Bot) StartStreaming() error {
	if s.stream != nil {
		panic("already have a stream")
	}

	connection, response, err := websocket.DefaultDialer.Dial(s.wsEndpoint, s.wsHeaders)
	if err != nil {
		s.logger.Fatal("dial:", err)
	}
	s.logger.Printf("received response from ws connection: status=%d headers=%v url=%s", response.StatusCode, response.Header, response.Request.URL)

	s.websocketConnection = connection

	s.session = joselito.NewSession(connection)
	s.stream = joselito.NewStream(s.Client, s.session)

	s.logger.Printf("joining talkgroups %v\n", s.talkgroups)
	err = s.session.GroupJoin(s.talkgroups)
	if err != nil {
		s.logger.Printf("can't join talkgroups %v: %v", s.talkgroups, err)
		s.Close()
		return err
	}

	return nil
}

func (s *Bot) StopStreaming() {
	s.logger.Printf("closing stream session")
	s.stream.Stop()
	s.Close()
	s.websocketConnection = nil
	s.session = nil
	s.stream = nil
}

func (s *Bot) Close() {
	if s.websocketConnection != nil {
		err := s.websocketConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			s.logger.Println("ws connection close:", err)
			return
		}

		s.websocketConnection.Close()
	}
}
