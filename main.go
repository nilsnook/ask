package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	LOGFILE := path.Join(os.TempDir(), "ask.log")
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// new ask
	a := newAsk(f)
	// setup config
	a.setupConfig()

	// set SSH_ASKPASS environment variable before
	// adding keys to ssh-agent
	// REQUIRED!!! because this program is what gets
	// us the password for a given key from password manager -
	// "pass" in my case.
	os.Setenv("SSH_ASKPASS", "asksshpass")

	// traverse through all keys in the config struct
	// and add them to ssh-agent
	var anyerr bool
	for _, key := range a.config.Keys {
		// resolve '~' or '$HOME' to absolute value
		key = a.resolveKey(key)
		// add keys if they exist
		if exists(key) {
			err := a.addKey(key)
			if err != nil {
				a.errorLog.Printf("Error adding key [%s]\n", key)
				a.errorLog.Fatalln(err)
				anyerr = true
			} else {
				a.infoLog.Printf("Added key [%s]\n", key)
			}
		} else {
			a.errorLog.Printf("Key [%s] does not exist. Ignoring...\n", key)
			anyerr = true
		}
	}

	fmt.Println("Added SSH Keys!")
	if anyerr {
		fmt.Println("(There were some errors. Check logs!)")
	}
}
