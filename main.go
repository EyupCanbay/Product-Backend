package main

import (
	"fmt"
	"tesodev_interview/configs"
)

func main() {
	fmt.Println("db bağlanma öncesi")
	configs.ConnectDB()
	fmt.Println("db bağlanma sonrası")

}
