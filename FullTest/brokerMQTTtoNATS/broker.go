package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")

func usage() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

func messageReceived(client *MQTT.Client, msg MQTT.Message) {
	topics := strings.Split(msg.Topic(), "/")
	msgFrom := topics[len(topics)-1]
	fmt.Print(msgFrom+": "+string(msg.Payload()), " to ")

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	//subj, msg := args[0], []byte(args[1])

	fmt.Print(msgFrom+": "+string(msg.Payload()), " to ", nc.ConnectedUrl())

	nc.Publish(msgFrom, msg.Payload())
	nc.Flush()
}

func main() {
	stdin := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	broker := flag.String("broker", "tcp://iot.eclipse.org:1883", "The MQTT broker to connect to")
	endpoint := flag.String("endpoint", "endpoint1", "The chat room to enter. default 'gochat'")
	name := flag.String("name", "user"+strconv.Itoa(rand.Intn(1000)), "Username to be displayed")
	flag.Parse()

	subTopic := strings.Join([]string{"/broker/", *endpoint, "/+"}, "")
	pubTopic := strings.Join([]string{"/broker/", *endpoint, "/", *name}, "")

	opts := MQTT.NewClientOptions().AddBroker(*broker).SetClientID(*name).SetCleanSession(true)

	opts.OnConnect = func(c *MQTT.Client) {
		if token := c.Subscribe(subTopic, 1, messageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Connected as %s to %s\n", *name, *broker)

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		if token := client.Publish(pubTopic, 1, false, message); token.Wait() && token.Error() != nil {
			fmt.Println("Failed to send message")
		}
	}
}
