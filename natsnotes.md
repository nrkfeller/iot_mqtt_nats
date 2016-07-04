# Nats + MQTT Architecture

Creating a IoT architecture that uses NATS to redirect information from the MQTT broker to a centralized system.

Sensors----Southbound(MQTT protocol)----> MQTT Broker ---->Northbound(NATS protocol)----> Centralized system

### Getting Started with NATS
Text based pub sub protocol. clients can communicate through gnatsd through regular TCP/IP socket. Uses subjects/topics and subtopics with widlcars, much the same was a MQTT.
```
telnet demo.nats.io 4222
```
* gnastd - go nats deamon, executable name for nats server
* server is highly configurable
* provided as docker image you can pull form docker hub
* Nats Top is a monitoring utility

### PubSub
Start server gnatsd.
Subscribe to a subject - go run sub-program.go planet.continent.country.town
```
go run nats-sub.go msg.test
```
```
go run nats-sub.go msg.*
```
Publish to a subject - go run pub-program.go planet.continent.* "lorem ipsum"
```
go run nats-pub.go msg.test "NATS MESSAGE"
```

### Queueing
Start server gnatsd
```
go run nats-sub.go foo my-queue
```
```
go run nats-sub.go foo
```
```
go run nats-pub.go foo "Hello NATS!"
```

### Monitoring with NATS
* General Stats: http://localhost:8222/varz
* Connection Details: http://localhost:8222/connz
* Subscriber Stats: http://localhost:8222/subscriptionsz
* Endpoint Routes: http://localhost:8222/routez
