package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/yeyus/gumble-joselito/pkg/bot"
	"github.com/yeyus/gumble-joselito/pkg/joselito"
	"github.com/yeyus/gumble/gumble"
)

type talkgroupList []*joselito.DMRID

func (i *talkgroupList) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *talkgroupList) Set(value string) error {
	tg, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		log.Panicf("can't parse talkgroup from value %s", value)
	}
	*i = append(*i, joselito.NewDMRID(uint(tg)))
	return nil
}

func main() {
	// Command line flags

	// Mumble section
	server := flag.String("server", "127.0.0.1:64738", "the server to connect to")
	username := flag.String("username", "", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", false, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")
	room := flag.String("room", "", "the Room path separated by commas where the streamer shall enter")

	// Websocket section
	endpoint := flag.String("endpoint", "", "the websocket endpoint to connect to")
	var talkgroups talkgroupList
	flag.Var(&talkgroups, "talkgroup", "list of comma separated talkgroup ids")
	userAgent := flag.String("useragent", "", "The user agent sent to the streaming server")
	origin := flag.String("origin", "", "The origin sent to the streaming server")

	flag.Parse()

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	httpHeaders := make(http.Header)
	if len(*origin) > 0 {
		httpHeaders.Add("Origin", *origin)
	}

	if len(*userAgent) > 0 {
		httpHeaders.Add("User-Agent", *userAgent)
	}

	if len(talkgroups) < 1 {
		log.Panicf("no talkgroup list was specified, halting")
	}

	// Create gumble client
	tlsConfig := &tls.Config{}
	if *insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	config := gumble.NewConfig()
	config.Username = *username
	config.Password = *password

	// Create websocket client
	log.Printf("connecting to %s", *endpoint)

	bot := bot.NewBot(
		*server,
		strings.Split(*room, ","),
		config,
		tlsConfig,
		*endpoint,
		httpHeaders,
		talkgroups,
	)
	defer bot.Close()

	bot.Connect()

	bot.WaitGroup.Wait()

	// for {
	// 	select {
	// 	case <-session.SessionEnd:
	// 		return
	// 	case <-interrupt:
	// 		log.Println("interrupt ctrl-c")

	// 		// Cleanly close the connection by sending a close message and then
	// 		// waiting (with timeout) for the server to close the connection.
	// 		bot.Close()
	// 		select {
	// 		case <-session.SessionEnd:
	// 		case <-time.After(time.Second):
	// 		}
	// 		return
	// 	}
	// }
}
