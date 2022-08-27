package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Jeffail/tunny"
	"github.com/xrjr/mcutils/pkg/ping"
)

func getIps(ips_file string) []string {
	file, err := os.Open(ips_file)
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

// CLI flags
var help = flag.Bool("help", false, "Show help")

var outfileFlag = "output.txt"
var ipsfileFlag = "ips.txt"
var numThreadsFlag = 1000
var portFlag = 25565
var startAtFlag = ""

func main() {
	// Binds the flags
	flag.StringVar(&outfileFlag, "o", "output.txt", "Output file (where all the data goes).")
	flag.StringVar(&ipsfileFlag, "i", "", "File containing all the IPv4s separated by linebreaks (input).")
	flag.StringVar(&startAtFlag, "r", "", "Start at this IPv4 (resume).")
	flag.IntVar(&numThreadsFlag, "n", 1000, "Number of scan workers (threads).")
	flag.IntVar(&portFlag, "p", 25565, "Server port to use.")

	// Parse the flags
	flag.Parse()

	// Validate dumb shit
	if ipsfileFlag == "" {
		fmt.Println("You must provide a input file!\nExample: `$ mcmap -i list-of-servers.txt`\nWrite `$ mcmap -h` for help.")
		os.Exit(1)
	}

	if numThreadsFlag <= 0 {
		fmt.Println(fmt.Sprintf("Number of threads must be >= 1, must provide more threads! %d <= 0!\nWrite `$ mcmap -h` for help.", numThreadsFlag))
		os.Exit(1)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Loop over ips and batch them into an array of 10 ips, then ping them in a goroutine
	allIps := getIps(ipsfileFlag)

	file, err := os.OpenFile(outfileFlag, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	pool := tunny.NewFunc(numThreadsFlag, func(data interface{}) interface{} {
		ip, ok := data.(string)
		if !ok {
			return nil
		}
		properties, ping, err := ping.Ping(ip, portFlag)
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
			"%s (ping: %d, players: %d, version: %s, motd: %s names: %s)\n",
			ip, ping, properties.Infos().Players.Online, properties.Infos().Version.Name,
			properties.Infos().Description, str,
		)

		file.WriteString(formattedData)
		fmt.Println(formattedData)

		return nil
	})

	defer pool.Close()
	shouldStart := startAtFlag == ""

	for _, ip := range allIps {
		if ip == startAtFlag {
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
}
