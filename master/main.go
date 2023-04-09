package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var workers = make([]string, 0)

var current_instances int = 0
var current_total int = 0
var current_good int = 0

func main() {
	app := fiber.New()

    app.Get("/status", func(c *fiber.Ctx) error {
        log.Println(workers)
        return c.SendStatus(fiber.StatusOK)
    })

    app.Get("/join", joinNetwork)
    app.Get("/remove", exitNetwork)
    app.Post("/collect", func (c *fiber.Ctx) error {
		var b map[string][]string
        
        if err := json.Unmarshal(c.Body(), &b); err != nil {
			return err
		}
        
        total, _ := strconv.Atoi(b["total"][0])
        valued, _ := strconv.Atoi(b["valued"][0])
        instances, _ := strconv.Atoi(b["instances"][0])

        current_total += total
        current_good += valued

        current_instances += 1

        if current_instances == instances {
            log.Printf("Finished: %f", float32(current_good)/float32(current_total))
            current_good = 0
            current_instances = 0
            current_total = 0
        }

        return c.SendStatus(fiber.StatusOK)
    })

    go func() {
        app.Listen(":3000")
    }()

    for {
        fmt.Print("Comanda: ")
        comanda := readFromStdinLine()
        fmt.Print("Date: ")
        data := readStringArrayFromStdinline();
        SendCommand(comanda, data)
    }    
}

func readFromStdinLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func readStringArrayFromStdinline() []string {
    data := make([]string, 0)
    for {
        cuvant := readFromStdinLine()
        if len(cuvant) == 0 {
            break
        }
        data = append(data, cuvant)
    }
    return data
} 