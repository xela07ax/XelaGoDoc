package main

import (
	"fmt"
	"net"
	"os"
)


func main() {
			availableInterfaces()
}

func availableInterfaces() {
	interfaces, err := net.Interfaces()

	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	fmt.Println("Available network interfaces on this machine : ")
	for _, i := range interfaces {
		fmt.Printf("Name : %v \n", i.Name)
	}
}
