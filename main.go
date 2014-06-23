package main

import (
	"crypto/rand"
	"fmt"
	"github.com/monochromegane/hoi/server"
	"os"
	"os/user"
	"path/filepath"
)

func main() {

	file := os.Args[1]

	// create hoi public directory
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	publicDir := filepath.Join(homeDir, ".hoi", "public")
	os.MkdirAll(publicDir, 0755)

	// create random directory
	random := randomString(32)
	randomDir := filepath.Join(publicDir, random)
	os.Mkdir(randomDir, 0755)

	// create symblic link
	os.Symlink(file, filepath.Join(randomDir, filepath.Base(file)))

	fmt.Println(filepath.Join(random, filepath.Base(file)))

	// start hoi server
	server.Start(publicDir)
}

func randomString(length int) string {
	alphanum := "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}