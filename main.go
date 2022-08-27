package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Jeffail/tunny"
	"github.com/xrjr/mcutils/pkg/ping"
)

func getIps() []string {
	// Read ips from file possible_minecraft_servers.txt
	file, err := os.Open("real_mc_servers.txt")
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
	allIps := getIps()

	var startAt string = "218.148.136.167"

	file, err := os.OpenFile("server-map12.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	pool := tunny.NewFunc(1000, func(data interface{}) interface{} {
		ip, ok := data.(string)
		if !ok {
			return nil
		}
		properties, ping, err := ping.Ping(ip, 25565)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		player_samples := properties.Infos().Players.Sample
		str := ""
		for _, sample := range player_samples {
			str += strings.ReplaceAll(sample.Name, "\n", " ") + ", "
		}

		formattedData := fmt.Sprintf(
			"%s (ping: %d, players: %d, version: %s, modt: %s names: %s)\n",
			ip, ping, properties.Infos().Players.Online, properties.Infos().Version.Name,
			properties.Infos().Description, str,
		)

		file.WriteString(formattedData)
		fmt.Println(formattedData)

		return nil
	})

	defer pool.Close()
	shouldStart := true

	for _, ip := range allIps {
		if ip == startAt {
			shouldStart = true
			continue
		}
		if !shouldStart {
			continue
		}

		go func(ip string) {
			pool.Process(ip)
		}(ip)
	}

	for pool.QueueLength() > 0 {
		time.Sleep(time.Second)
	}

	// Write servers with players and is 1.19.X to file servers_with_players_and_is_1.19.X.txt
}
