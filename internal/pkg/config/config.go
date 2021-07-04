package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/mineway/worker/internal/pkg/rig"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	// Embed
	BuildName 		string		`json:"build_name"`
	Version 		string 		`json:"version"`
	// Process Chan
	ApiChan 		chan string
	// Detect
	HomeDir 		string
	ExcavatorDir 	string
	OS 				string
	Arch			string
	RigData 		*rig.Data
	// User Current Setting
	userConfig
	// Utils
	CreatedAt   	time.Time
}

type userConfig struct {
	RigName 		string		`json:"rig_name"`
	Miner 			string		`json:"miner"`
	Algo			string		`json:"algo"`
	WebInterfaceURL	string		`json:"web_interface_url"`
	ApiPort 		string		`json:"api_port"`
}

func New() (*Config, error) {
	rd, err := rig.New()
	return &Config{
		RigData: rd,
	}, err
}

// Init allows to search exist configuration or create it
func (c *Config) Init() (err error) {
	dir, err := homedir.Dir()
	if err != nil {
		return
	}

	c.HomeDir, err = homedir.Expand(dir)
	if err != nil {
		return
	}

	c.ExcavatorDir = filepath.Join(c.HomeDir, "." + c.BuildName)

	_, err = os.Stat(c.ExcavatorDir)
	if err != nil {
		// Try to create it if not exist
		if err = c.create(); err != nil {
			return
		}
	}

	return c.read()
}

// Save allows to save current configuration into json file
func (c *Config) Save() error {
	uc := userConfig{
		RigName: c.RigName,
		Miner: c.Miner,
		Algo: c.Algo,
	}
	file, err := json.MarshalIndent(uc, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		filepath.Join(c.ExcavatorDir, "config.json"),
		file,
		0600,
		)
}

func (c *Config) StopAPI() {
	c.ApiChan <- "stop"
}

// create allows to create config directory into home's directory
func (c *Config) create() (err error) {
	d, err := os.Stat(c.ExcavatorDir)
	if err == nil && !d.IsDir() {
		return fmt.Errorf("%s already exist but isn't a directory", c.ExcavatorDir)
	}

	for {
		c.RigName = c.formValue("Enter your rig name")
		c.WebInterfaceURL = c.formValue("Enter your web interface url")

		restart := c.formValid()
		if !restart {
			screen.Clear()
			break
		}
	}

	err = os.Mkdir(c.ExcavatorDir, 0700)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(c.ExcavatorDir, "config.json"))
	if err != nil {
		return err
	}

	//var rgx = regexp.MustCompile(`(?m)("[^"]+"|[^\s"]+)`)
	return c.Save()
}

// Import config from config directory into memories
func (c *Config) read() error {
	data, err := ioutil.ReadFile(filepath.Join(c.ExcavatorDir, "config.json"))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c)
}

func (c *Config) welcomeHeader() {
	screen.Clear()
	color.Magenta("Welcome to %s's install\n", c.BuildName)
	color.Yellow("Set-up your worker here, your configuration file will be here if you need to edit it :")
	color.Yellow("%s", c.ExcavatorDir)
	color.Yellow("------------------------------")
}

func (c *Config) formValue(question string) string {
	var str string

	for {
		c.welcomeHeader()

		color.Cyan("%s :\n", question)

		reader := bufio.NewReader(os.Stdin)

		str, _ = reader.ReadString('\n')
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, "\"", "", -1)
		str = strings.TrimSpace(str)

		if len(str) != 0 {
			break
		}
	}

	return str
}

func (c *Config) formValid() bool {
	c.welcomeHeader()
	color.Yellow("Your settings are :")
	color.Yellow("Rig name => %s", color.GreenString(c.RigName))
	color.Yellow("Web interface url => %s", color.GreenString(c.WebInterfaceURL))
	color.Yellow("------------------------------")
	color.Cyan("Do you want restart set-up ? [y/N]\n")

	reader := bufio.NewReader(os.Stdin)

	str, _ := reader.ReadString('\n')
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\"", "", -1)

	if strings.ToLower(strings.TrimSpace(str)) == "y" {
		return true
	}
	return false
}