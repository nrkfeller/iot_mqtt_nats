package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

func main() {
	//stdin := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	broker := flag.String("broker", "tcp://iot.eclipse.org:1883", "The MQTT broker to connect to")
	endpoint := flag.String("endpoint", "endpoint1", "The chat room to enter. default 'gochat'")
	name := flag.String("name", "brokerrouter", "Username to be displayed")
	flag.Parse()

	pubTopic := strings.Join([]string{"/broker/", *endpoint, "/", *name}, "")

	opts := MQTT.NewClientOptions().AddBroker(*broker).SetClientID(*name).SetCleanSession(true)

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Connected as %s to %s\n", *name, *broker)

	message := 1

	go func() {
		for {
			for {
				// message, err := stdin.ReadString('\n')
				// if err == io.EOF {
				// 	os.Exit(0)
				// }
				if token := client.Publish(pubTopic, 1, false, fmt.Sprintf("%d", message)); token.Wait() && token.Error() != nil {
					fmt.Println("Failed to send message")
				}
				message++
			}
		}
	}()

	time.Sleep(time.Second)
	fmt.Println(message)
}
