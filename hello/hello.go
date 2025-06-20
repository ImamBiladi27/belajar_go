package main

import (
    "fmt"
	 "log"
    "example.com/greetings"
)
func main() {
    message, err := greetings.Hello("Imam")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(message)
}
// func main() {
//     message := greetings.Hello("Imam")
//     fmt.Println(message)
// }
