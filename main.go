package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/xrjr/mcutils/pkg/ping"
)

func getIps() []string {
	// Read ips from file possible_minecraft_servers.txt
	file, err := os.Open("possible_minecraft_servers.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var ips []string
	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}
	return ips
}

var allNodesWaitGroup sync.WaitGroup

func main() {
	// // Loop over ips and ping them
	// for _, ip := range getIps() {
	// 	properties, _, err := ping.Ping(ip, 25565)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(properties.Infos().Players)
	// }

	// Loop over ips and ping them in a goroutine
	// for _, ip := range getIps() {
	// 	go func(ip string) {
	// 		properties, _, err := ping.Ping(ip, 25565)
	// 		if err != nil {
	// 			return
	// 		}
	// 		fmt.Println(properties.Infos().Players)
	// 	}(ip)
	// }

	// Loop over ips and batch them into an array of 10 ips, then ping them in a goroutine
	perRoutine := 20
	allIps := getIps()

	file, err := os.Create("5-servers_with_players_and_is_1.19.X.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for i := 0; i < len(allIps); i += perRoutine {
		ips := allIps[i : i+perRoutine]
		allNodesWaitGroup.Add(1)
		go func(ips []string) {
			defer allNodesWaitGroup.Done()
			for _, ip := range ips {
				properties, _, err := ping.Ping(ip, 25565)
				if err != nil {
					fmt.Println(err)
					return
				}
				if strings.Contains(properties.Infos().Version.Name, "1.19") {
					if properties.Infos().Players.Online > 0 {
						file.WriteString(ip + "\n")
					}
				}
				fmt.Println("IP:", ip, "Players:", properties.Infos().Players)
			}
		}(ips)
	}

	// Write servers with players and is 1.19.X to file servers_with_players_and_is_1.19.X.txt

	allNodesWaitGroup.Wait()
}
