package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Helper: create HTTP client for Unix socket
func newUnixClient(socketPath string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
		Timeout: 5 * time.Second,
	}
}

// Send request to the server
func sendRequest(client *http.Client, path string) (string, error) {
	resp, err := client.Get("http://unix" + path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	return string(buf[:n]), nil
}

func main() {
	socketPath := "/tmp/rackkv.sock"
	client := newUnixClient(socketPath)

	fmt.Println("Welcome to rackKV CLI")
	fmt.Println("Type 'help' for commands")

	scanner := bufio.NewScanner(os.Stdin)
	banner := ` ________  ________  ________  ___  __    ___  __    ___      ___ 
|\   __  \|\   __  \|\   ____\|\  \|\  \ |\  \|\  \ |\  \    /  /|
\ \  \|\  \ \  \|\  \ \  \___|\ \  \/  /|\ \  \/  /|\ \  \  /  / /
 \ \   _  _\ \   __  \ \  \    \ \   ___  \ \   ___  \ \  \/  / / 
  \ \  \\  \\ \  \ \  \ \  \____\ \  \\ \  \ \  \\ \  \ \    / /  
   \ \__\\ _\\ \__\ \__\ \_______\ \__\\ \__\ \__\\ \__\ \__/ /   
    \|__|\|__|\|__|\|__|\|_______|\|__| \|__|\|__| \|__|\|__|/    
                                                                  
                                                                  
                                                                  `
	fmt.Println(banner)

	for {
		fmt.Print("rackKV> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "HELP":
			fmt.Println("Commands:")
			fmt.Println("  OPEN -rw -sync           Open database with optional flags")
			fmt.Println("  PUT <key> <value>        Put key-value")
			fmt.Println("  GET <key>                Get value by key")
			fmt.Println("  DELETE <key>             Delete key")
			fmt.Println("  EXIT                     Quit CLI")

		case "EXIT":
			fmt.Println("Bye!")
			return

		case "OPEN":
			rw := "false"
			sync := "false"
			for _, arg := range parts[1:] {
				switch strings.ToLower(arg) {
				case "-rw":
					rw = "true"
				case "-sync":
					sync = "true"
				}
			}
			res, err := sendRequest(client, fmt.Sprintf("/open?rw=%s&syn=%s", rw, sync))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(res)
			}

		case "PUT":
			if len(parts) < 3 {
				fmt.Println("Usage: PUT <key> <value>")
				continue
			}
			key := parts[1]
			val := strings.Join(parts[2:], " ")
			res, err := sendRequest(client, fmt.Sprintf("/put?key=%s&value=%s", key, val))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(res)
			}

		case "GET":
			if len(parts) != 2 {
				fmt.Println("Usage: GET <key>")
				continue
			}
			key := parts[1]
			res, err := sendRequest(client, fmt.Sprintf("/get?key=%s", key))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(res)
			}

		case "DELETE":
			if len(parts) != 2 {
				fmt.Println("Usage: DELETE <key>")
				continue
			}
			key := parts[1]
			res, err := sendRequest(client, fmt.Sprintf("/delete?key=%s", key))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(res)
			}

		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
