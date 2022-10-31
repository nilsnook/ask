package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type ask struct {
	infoLog   *log.Logger
	errorLog  *log.Logger
	homeDir   string
	configDir string
	config    config
}

func newAsk(f *os.File) *ask {
	infolog := log.New(f, "INFO\t", log.LstdFlags)
	errorlog := log.New(f, "ERROR\t", log.LstdFlags|log.Lshortfile)

	userhomedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	userconfigdir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	askconfigdir := path.Join(userconfigdir, "ask/")

	return &ask{
		infoLog:   infolog,
		errorLog:  errorlog,
		homeDir:   userhomedir,
		configDir: askconfigdir,
	}
}

func (a *ask) setupConfig() {
	// logging empty line before each new config setup
	a.infoLog.Println()
	// init config
	c := config{}
	// set defaults
	c.setDefaults()

	// read configuration from config file
	if err := c.readConfigFileIn(a.configDir); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			a.errorLog.Printf("Config file not found (in %s)... Using defaults!\n", a.configDir)
			// config file not found, write one with defaults
			if err = c.writeConfigFileIn(a.configDir); err != nil {
				a.errorLog.Printf("Error writing config file (in %s)\n", a.config)
				a.errorLog.Fatalln(err)
			}
		} else {
			// config file was found but another error was reported
			a.errorLog.Println("Config file found, but another error was reported")
			a.errorLog.Fatalln(err)
		}
	}

	// unmarshall config into struct
	err := viper.Unmarshal(&c)
	if err != nil {
		a.errorLog.Println("Unable to decode config file")
		a.errorLog.Fatalln(err)
	}

	// set config
	a.config = c
}

func (a *ask) resolveKey(key string) string {
	if strings.HasPrefix(key, "~/") {
		key = strings.TrimLeft(key, "~/")
		key = path.Join(a.homeDir, key)
	} else if strings.HasPrefix(key, "$HOME/") {
		key = strings.TrimLeft(key, "$HOME/")
		key = path.Join(a.homeDir, key)
	}
	return key
}

func (a *ask) addKey(key string) error {
	sshaddCmd := exec.Command("ssh-add", key)
	sshaddCmd.Stdout = a.infoLog.Writer()
	sshaddCmd.Stderr = a.errorLog.Writer()
	err := sshaddCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
