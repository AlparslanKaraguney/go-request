# go-request

Simple HTTP Request library for golang.

Usage:
```
package main

import (
	"fmt"
	"github.com/AlparslanKaraguney/go-request"
)

func main() {

	client := request.NewHttpRequestClient()

	response, err := client.Get("https://www.google.com", nil)

	if err != nil {
		fmt.Println("Some error occurred:", err.Error())
		return
	}

	fmt.Print(response.StatusCode)

}
```