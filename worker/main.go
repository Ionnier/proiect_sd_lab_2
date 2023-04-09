package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"common"

	"github.com/gofiber/fiber/v2"
)
var empty_map = make(map [string]string, 0)

func main() {
	master_url := os.Getenv("MASTER_HOSTNAME")
	my_hostname := os.Getenv("MY_HOSTNAME")
	shuffler := os.Getenv("SHUFFLER")


	if (len(master_url) == 0 || len(my_hostname) == 0) {
		master_url = "http://localhost:3000"
		my_hostname = "http://localhost:4000"
		shuffler = "http://localhost:5000"
		log.Printf("master_url= %s, my_hostname=%s", master_url, my_hostname)
	}

	splits := strings.Split(my_hostname, ":")
	port := splits[len(splits)-1]

	log.Printf(
		"Initialized with:\n master_url=%s\n my_hostname=%s\n port=%s\n",
		master_url,
		my_hostname,
		port,
	)
	
	for {
		my_map := make(map[string] string, 0)
		my_map["serverip"] = my_hostname
		status, data := common.SendRequest(
			fmt.Sprintf("%s/join", master_url),
			empty_map,
			fiber.MethodGet,
			my_map,
		)
		if status == fiber.StatusOK {
			break;
		}

		log.Printf("Received %v %s... waiting to retry", status, string(data))
		time.Sleep(15000)
	}

	log.Printf("Worker %s initialised, starting server", my_hostname)

	app := fiber.New()

    app.Get("/status", func(c *fiber.Ctx) error {
		log.Println("Got pinged...")
        return c.SendStatus(fiber.StatusOK)
    })

	app.Post("/reduce", func(c *fiber.Ctx) error {
		var b map[string][]string

		if err := json.Unmarshal(c.Body(), &b); err != nil {
			return err
		}

		date := b["data"]

		reduceMap := make(map [string] int, 0)
		for _, v := range date {
			aux := strings.Split(v, ":")
			value := aux[0]
			occurances := aux[1]
			val, _ := strconv.Atoi(occurances)
			
			if el, ok := reduceMap[value]; ok {
				reduceMap[value] = el + val 
			} else {
				reduceMap[value] = val
			}

		}

		total := 0
		exista := 0
		comanda := b["command"][0]
		// instances := b["instances"][0]
		for key, index := range reduceMap {
			total += index
			comanda = "ex4"
			if comanda == "ex4" && checkVowels(key) {
				exista += index
			} else if comanda == "ex5" {
				exista += index
			}
		}
		to_send := make(map[string][]string, 0)
		to_send["total"] = []string{fmt.Sprintf("%d", total)}
		to_send["valued"] = []string{fmt.Sprintf("%d", exista)}
		to_send["instances"] = b["instances"]

		common.SendRequest(
			fmt.Sprintf("%s/collect", master_url),
			to_send,
			fiber.MethodPost,
			make(map[string]string,0),
		)

		log.Printf("Am redus din %v in %d si %d", reduceMap, total, exista)

        return c.SendStatus(fiber.StatusOK)
    })

	app.Post("/map", func(c *fiber.Ctx) error {
		var b map[string][]string

		if err := json.Unmarshal(c.Body(), &b); err != nil {
			return err
		}
		comanda := b["command"][0]
		// instances := b["instances"][0]
		date := b["data"]

		log.Printf("%s face map cu comanda %s pe %v\n", my_hostname, comanda, date)

		mapa := make(map[string] int, 0)
		for _, val := range date {
			if value, ok := mapa[val]; ok {
				mapa[val] = value + 1
			} else {
				mapa[val] = 1
			}
		}

		log.Printf("%s a mapat %v", my_hostname, mapa)

		to_send := make([]string, 0)

		for i, v := range mapa {
			to_send = append(to_send, fmt.Sprintf("%s:%d", i, v))
		}
		var data_send = make(map[string][]string, 0)
		data_send["instances"]= b["instances"]
		data_send["data"] = to_send
		data_send["myip"] = []string{my_hostname}
		data_send["command"] = b["command"]
		common.SendRequest(
			fmt.Sprintf("%s/shuffle", shuffler),
			data_send,
			fiber.MethodPost,
			make(map[string] string, 0),
		)

        return c.SendStatus(fiber.StatusOK)
    })

	app.Listen(fmt.Sprintf(":%s", port))
}

func checkVowels(cuvant string) bool {
	vowels := []string{"a", "e", "i", "o", "u"}
	if isElementExist(vowels, cuvant[0]) && isElementExist(vowels, cuvant[len(cuvant)-1]) {
		return true
	}
	return false
}


func isElementExist(s []string, str byte) bool {
	for _, v := range s {
	  if v == string(str) {
		return true
	  }
	}
	return false
  }