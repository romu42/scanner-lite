package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Read the config file of zones to scan.
func ReadConf() map[string]*ZoneP {
	zones := make(map[string]*ZoneP)

	file, err := os.Open("zones2scan.yml")
	log.Printf("Reading %s for zones to scan\n", file.Name())
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#") {
			continue
		}

		// For each line in the list create Zone Parent struct
		// ex: msat1.catch22.se.:ns1.catch22.se.:13.48.238.90:53:hmac-sha256:musiclab.parent:4ytnbnbTtA+w19eZjK6bjw/VB9SH8q/5eQKvf9BlAf8=
		parts := strings.Split(line, ":")
		z := &ZoneP{
			hostname: parts[1],
			ip:       parts[2],
			port:     parts[3],
			hmac:     parts[4],
			keyname:  parts[5],
			secret:   parts[6],
		}
		zones[parts[0]] = z
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return zones
}
