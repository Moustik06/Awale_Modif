package Game

import (
	"fmt"
	"os"
	"os/signal"
	. "projet-ai/Types"
	. "projet-ai/utils"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

/*
BrokerAddress -> le serveur MQTT utilisé
myName -> le nom du joueur
opponentName -> le nom de l'adversaire

Ce code pour jouer en MQTT ne fonctionne qu'en Go, car j'utilise des channels pour communiquer entre des goroutines.
Les goroutines sont des threads légers qui permettent de faire de la concurrence en Go.
*/
const brokerAddress string = "test.mosquitto.org:1883"
const myName string = "Quentin"
const opponentName string = "Julien"

// messageHandler est une fonction qui traite les messages reçus par le client MQTT.
// Elle prend en entrée un client MQTT et un message MQTT.
// La fonction extrait le mouvement reçu du message et l'envoie sur le canal ReceivedMove.
var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	move := string(msg.Payload())
	var color MoveColor
	fmt.Printf("\nMouvement reçu : %s\n", move)
	index, line := SplitAfterLastInt(move)
	switch line {
	case "B":
		color = B
	case "R":
		color = R
	case "TB":
		color = TB
	default:
		color = TR
	}
	final := Move{HoleIndex: index - 1, MoveColor: color}
	ReceivedMove <- final
}

// sendMessage est une fonction qui envoie un message MQTT à l'adversaire.
// Elle prend en entrée un client MQTT et un message.
// La fonction publie le message sur le topic "awale/"+opponentName.
func sendMessage(client MQTT.Client, message string) {
	fmt.Println("Mouvement envoyé : ", message)
	token := client.Publish("awale/"+opponentName, 0, false, message)
	token.Wait()
}

// InitMqtt initialise un client MQTT et s'abonne à un topic pour recevoir des messages.
// Elle prend en entrée un canal de type Move et envoie des messages sur le topic abonné.
// La fonction bloque jusqu'à ce qu'un message soit reçu sur le canal.
func InitMqtt(ch chan Move) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID(myName)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	client.Subscribe("awale/"+myName, 0, messageHandler)
	defer client.Disconnect(250)
	fmt.Println("Connecté au broker MQTT")
	for {
		select {
		case move := <-ch:
			var result = ""
			println("Player choose move: ", move.HoleIndex+1, move.MoveColor)
			result += strconv.Itoa(move.HoleIndex + 1)
			switch move.MoveColor {
			case B:
				result += "B"
			case R:
				result += "R"
			case TB:
				result += "TB"
			default:
				result += "TR"
			}

			sendMessage(client, result)
		}
	}

}
