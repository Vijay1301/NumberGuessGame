package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var secretNumber int

func GenerateNumber() {
	rand.Seed(time.Now().UnixNano())
	secretNumber = rand.Intn(1000) + 1
}

func Socket(app *fiber.App) {
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			guess, err := strconv.Atoi(string(msg))
			if err != nil {
				c.WriteMessage(messageType, []byte("Invalid number. Try again!"))
				continue
			}

			response := ""
			if guess < secretNumber {
				response = "Too low! Try again."
			} else if guess > secretNumber {
				response = "Too high! Try again."
			} else {
				response = fmt.Sprintf("Congratulations! You guessed it right. The number was %d", secretNumber)
				secretNumber = rand.Intn(1000) + 1
			}

			c.WriteMessage(messageType, []byte(response))
		}
	}))
}

func main() {

	GenerateNumber()

	app := fiber.New()

	Socket(app)

	log.Fatal(app.Listen(":3000"))
}
