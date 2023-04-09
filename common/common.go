package common

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type test struct {
	url string
}

func (m test) Write(p []byte) (n int, err error) {
	a := fiber.AcquireAgent()
	data := make(map[string]string)
	data["message"] = string(p)
	a.JSON(data)
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(m.url)
	if err := a.Parse(); err != nil {
		panic(err)
	}
	if _, _, err := a.Bytes(); err != nil {
		fmt.Print(err)
	}
	return fmt.Print(string(p))
}

func GetHTTPLogger(url string) test {
	var asd test
	asd.url = url
	return asd
}

func SendRequest(url string, data interface{}, method string, headers map[string]string) (int, []byte){
	log.Printf("Start %s request to %s with data=%v and headers %v", method, url, data, headers)
	a := fiber.AcquireAgent()
	a.JSON(data)
	req := a.Request()
	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if err := a.Parse(); err != nil {
		panic(err)
	}
	if status, data, err := a.Bytes(); err != nil {
		log.Println(err)
		return status, data
	} else {
		return status, data
	}
}

