package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	//"copy-udp/headers"
)

func get_address(passed_addr string) string {

	if passed_addr == "" {
		passed_addr = ":3000"
	}

	split_arr := strings.Split(passed_addr, ":")

	split_addr := split_arr[0]

	fmt.Printf("%s\n", split_addr)

	ip_arr, err := net.LookupAddr(split_addr)

	ip_addr := ""

	fmt.Print(ip_arr)

	// if len(ip_arr) > 1 {
	//   ip_addr = ip_arr[1]
	// } else {
	//   ip_addr = ip_arr[0]
	// }

	if err != nil {
		log.Fatalf("Unable to find the given address")
	}

	if len(split_arr) > 1 {
		return fmt.Sprintf("%s:%s", ip_addr, split_arr[1])
	}

	return ip_addr
}

func main() {
	passed_addr := ""

	if len(os.Args) > 1 {
		passed_addr = os.Args[1]
	}

	addr := get_address(passed_addr)

	fmt.Printf("%s\n", addr)

	conn, err := net.Dial("ip4:udp", addr)

	fmt.Printf("%s", conn.LocalAddr())

	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
