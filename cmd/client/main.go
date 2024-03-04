package main

import (
	"flag"
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		address  = flag.String("address", "127.0.0.1:2222", "ssh target address")
		user     = flag.String("username", "test", "ssh user")
		password = flag.String("password", "test", "ssh password")
	)
	flag.Parse()

	client, err := ssh.Dial("tcp", *address, &ssh.ClientConfig{
		User:            *user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
	})
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	output, err := session.Output("uname -a")
	if err != nil {
		return err
	}

	log.Print(string(output))
	return nil
}
