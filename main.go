package main

import "github.com/kanyarats/finalexam/customers"

func main() {
	r := customers.SetupHandler()
	r.Run(":2019")
}
