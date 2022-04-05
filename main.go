package main

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/muesli/coral"
)

var (
	host     string
	password string
	port     uint32
	version  string

	rootCmd = &coral.Command{
		Use:   "obs-cli",
		Short: "obs-cli is a command-line remote control for OBS",
	}

	stdinCmd = &coral.Command{
		Use:   "stdin",
		Short: "Reads lines on stdin and executes each as a command",
		RunE: func(cmd *coral.Command, args []string) error {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				os.Args = append([]string{os.Args[0]},strings.Split(scanner.Text(), " ")...)
				rootCmd.Execute()
			}
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}

	client *goobs.Client
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if client != nil {
		_ = client.Disconnect()
	}
}

func init() {
	coral.OnInitialize(connectOBS)
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "host to connect to")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "password for connection")
	rootCmd.PersistentFlags().Uint32VarP(&port, "port", "p", 4444, "port to connect to")
	rootCmd.AddCommand(stdinCmd)
}

func getUserAgent() string {
	userAgent := "obs-cli"
	if version != "" {
		userAgent += "/" + version
	}
	return userAgent
}

func connectOBS() {
	var err error
	if client == nil {
		client, err = goobs.New(
			host+fmt.Sprintf(":%d", port),
			goobs.WithPassword(password),
			goobs.WithRequestHeader(http.Header{"User-Agent": []string{getUserAgent()}}),
		)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}
