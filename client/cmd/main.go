package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// request body (payload)
	requestBody := strings.NewReader(`
		{
			"msg_id": 1,
			"sender": "Tom",
			"msg": "Hello"
		}
	`)

	// post some data
	res, err := http.Post(
		"http://localhost:8080/greeter",
		"application/json; charset=UTF-8",
		requestBody,
	)

	// check for response error
	if err != nil {
		log.Fatal(err)
	}

	// read response data
	data, _ := ioutil.ReadAll(res.Body)

	// close response body
	res.Body.Close()

	// print request `Content-Type` header
	requestContentType := res.Request.Header.Get("Content-Type")
	fmt.Println("Request content-type:", requestContentType)

	// print response body
	fmt.Printf("%s\n", data)

}
