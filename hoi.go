package hoi

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Hoi struct {
	publicDir string
	config    Config
}

func NewHoi() *Hoi {
	// create public dir
	publicDir := publicDir()
	os.MkdirAll(publicDir, 0755)

	return &Hoi{
		publicDir: publicDir,
		config:    Load(configPath()),
	}
}

func (h Hoi) TestFile(file string) (string, error) {
	var (
		path string
		err  error
	)
	// absolutize
	path, err = filepath.Abs(file)
	if err != nil {
		return path, err
	}
	// check existence
	_, err = os.Stat(path)
	if err != nil {
		return path, err
	}
	return path, nil
}

func (h Hoi) MakePublic(file string) string {
	linked := h.makePublic(file)
	return linked
}

func (h Hoi) makePublic(src string) string {
	// create random directory
	random := h.createRandomDir()

	// make public by symblic link
	file := filepath.Base(src)
	os.Symlink(src, filepath.Join(h.publicDir, random, file))

	return filepath.Join(random, file)
}

func (h Hoi) MakeMessage(msgs []string) string {
	message := h.makeMessage(msgs)
	return message
}

func (h Hoi) makeMessage(msgs []string) string {
	// create random directory
	random := h.createRandomDir()

	// make public by message file
	file, err := os.Create(filepath.Join(h.publicDir, random, "message.txt"))
	if err != nil {
		return ""
	}
	defer file.Close()

	file.WriteString(strings.Join(msgs, " "))
	return filepath.Join(random, "message.txt")
}

func (h Hoi) createRandomDir() string {
	random := randomString(32)
	os.Mkdir(filepath.Join(h.publicDir, random), 0755)
	return random
}

func (h Hoi) Server() *HoiServer {
	return &HoiServer{
		DocumentRoot: h.publicDir,
		Port:         h.config.Port,
	}
}

func (h Hoi) Clear() {
	contents, _ := ioutil.ReadDir(h.publicDir)
	for _, c := range contents {
		os.RemoveAll(filepath.Join(h.publicDir, c.Name()))
	}
}

func (h Hoi) ToUrl(path string) string {
	server := h.Server()
	return fmt.Sprintf("%s/%s", server.Url(), path)
}

func (h Hoi) Notify(to, message string) string {
	n := NewNotifier(h.config.Notification)
	if n == nil {
		return ""
	}
	if err := n.Notify(to, message); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Message sent successfully to %s\n", to)
}

func publicDir() string {
	return filepath.Join(homeDir(), ".hoi", "temp_public")
}

func configPath() string {
	return filepath.Join(homeDir(), ".hoi", "conf.json")
}

func homeDir() string {
	usr, err := user.Current()
	var homeDir string
	if err == nil {
		homeDir = usr.HomeDir
	} else {
		// Maybe it's cross compilation without cgo support. (darwin, unix)
		homeDir = os.Getenv("HOME")
	}
	return homeDir
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
