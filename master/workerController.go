package main

import (
	"fmt"
	"log"

	"common"

	"github.com/gofiber/fiber/v2"
)

func joinNetwork(c *fiber.Ctx) error {
	server_ip := c.Get("serverip")
	if len(server_ip) == 0{
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, v := range workers {
		if v == server_ip {
			return c.SendStatus(fiber.StatusOK)
		}
	}

	workers = append(workers, server_ip)
	log.Println(workers)
	return c.SendStatus(fiber.StatusOK)
}

func exitNetwork(c *fiber.Ctx) error {
	server_ip := c.Get("serverip")
	if len(server_ip) == 0{
		return c.SendStatus(fiber.StatusBadRequest)
	}

	i := 0
	for _, v := range workers {
		if v == server_ip {
			workers = append(workers[:i], workers[i:]...)
			return c.SendStatus(fiber.StatusOK)
		}
		i++
	}

	return c.SendStatus(fiber.StatusBadRequest)
}

func SendCommand(command string, data []string) error {
	currentworkers := make([]string, 0 )

	for _, v := range(workers) {
		status, _ := common.SendRequest(fmt.Sprintf("%s/status", v), make(map[string]string, 0), fiber.MethodGet, make(map[string]string, 0))
		if (status == fiber.StatusOK) {
			currentworkers = append(currentworkers, v)
		}
	}

	log.Printf("%v", currentworkers)

	splits := SplitSlice(data, len(currentworkers))

	for i:=0; i<len(currentworkers); i++ {
		send := make(map[string][]string, 0)
		comman := make([]string, 0)
		comman = append(comman, command)
		send["command"] = comman

		instance := make([]string, 0)
		instance = append(instance, fmt.Sprintf("%d", len(currentworkers)))

		send["instances"] = instance
		send["data"] = splits[i]

		status, _ := common.SendRequest(fmt.Sprintf("%s/map", currentworkers[i]), send, fiber.MethodPost, make(map[string]string, 0))
		
		if (status != fiber.StatusOK) {
			log.Fatalf("Failed %s", currentworkers[i])
		}
		log.Printf("%s data %v", currentworkers[i], splits[i])
	}
	
	return nil
}

func SplitSlice(array []string, numberOfChunks int) [][]string {
	if len(array) == 0 {
		return nil
	}
	if numberOfChunks <= 0 {
		return nil
	}

	if numberOfChunks == 1 {
		return [][]string{array}
	}

	result := make([][]string, numberOfChunks)

	// we have more splits than elements in the input array.
	if numberOfChunks > len(array) {
		for i := 0; i < len(array); i++ {
			result[i] = []string{array[i]}
		}
		return result
	}

	for i := 0; i < numberOfChunks; i++ {

		min := (i * len(array) / numberOfChunks)
		max := ((i + 1) * len(array)) / numberOfChunks

		result[i] = array[min:max]

	}

	return result
}