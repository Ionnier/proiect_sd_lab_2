package main

import (
	"common"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var max_instances = 0
var workers = make([]string, 0)
var elements = make([]element, 0)

type element struct{
	Count string
	Data string
}

// Define a collection of that type
type Elements []element
 
// implement the functions from the sort.Interface
func (d Elements) Len() int{ 
    return len(d) 
}
 
func (d Elements) Less(i, j int) bool{ 
    return d[i].Data < d[j].Data 
}
 
func (d Elements) Swap(i, j int){ 
    d[i], d[j] = d[j], d[i] 
}

func main() {
	app := fiber.New()

	app.Post("/shuffle", func(c *fiber.Ctx) error {
		log.Println(max_instances)
		var b map[string][]string

		if err := json.Unmarshal(c.Body(), &b); err != nil {
			return err
		}
		// comanda := b["command"][0]
		instances := b["instances"][0]
		date := b["data"]
		workers = append(workers, b["myip"][0])

		for _, cuvinte := range date {
			process := strings.Split(cuvinte, ":")
			val := process[0]
			instances := process[1]
			elements = append(elements, element{
				Count: instances,
				Data: val,
			})
		}
		
		max_instances = max_instances + 1;

		if s, err := strconv.Atoi(instances); err == nil {
			if (s == max_instances) {
				sort.Sort(Elements(elements))
				for _, el := range elements {
					fmt.Println(el.Count, el.Data)
				}
				result := make([][]element, max_instances)

				curr_chunk := 0

				min_elements := len(elements) / max_instances

				for i:=0; i<len(elements); i++ {
					fmt.Println(elements[i])
					result[curr_chunk] = append(result[curr_chunk], elements[i])
					shouldSkip := false
					for i+1 < len(elements) && elements[i+1].Data == result[curr_chunk][len(result[curr_chunk])-1].Data {
						result[curr_chunk] = append(result[curr_chunk], elements[i+1])
						i = i+1
						if i >= len(elements) {
							shouldSkip = true
							break
						}
					}
					if shouldSkip {
						break
					}
					if len(result[curr_chunk]) > min_elements {
						curr_chunk += 1
					}
				}

				for i, v := range result {
					if (len(v) == 0 ){
						continue
					}
					log.Printf("send to reducer %s %v", workers[i], v)
					to_send := make(map[string][]string)
					to_send["command"] = b["command"]
					to_send["instances"] = b["instances"]
					data := make([]string, 0)
					for _, d := range v {
						data = append(data, fmt.Sprintf("%s:%s", d.Data, d.Count))
					}
					to_send["data"] = data

					common.SendRequest(
						fmt.Sprintf("%s/reduce", workers[i]),
						to_send,
						fiber.MethodPost,
						make(map[string]string, 0),
					)
				}
				max_instances = 0
				elements = make([]element, 0)


			}
		}

        return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":5000")
}

func SplitSlice(array []element, numberOfChunks int) [][]element {
	if len(array) == 0 {
		return nil
	}
	if numberOfChunks <= 0 {
		return nil
	}

	if numberOfChunks == 1 {
		return [][]element{array}
	}

	result := make([][]element, numberOfChunks)

	// we have more splits than elements in the input array.
	if numberOfChunks > len(array) {
		for i := 0; i < len(array); i++ {
			result[i] = []element{array[i]}
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