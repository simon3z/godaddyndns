package main

// cspell:ignore godaddyndns

import (
	"flag"
	"log"
	"time"

	"github.com/simon3z/godaddyndns"
	"github.com/simon3z/godaddyndns/cmd"
)

var cmdFlags = struct {
	ConfigFilePath string
	CheckInterval  int64
}{}

var cmdConfig *cmd.Config

func init() {
	flag.StringVar(&cmdFlags.ConfigFilePath, "c", "", "configuration file path")
	flag.Int64Var(&cmdFlags.CheckInterval, "i", 300, "check interval in seconds")
}

func CheckAndUpdate() error {
	extIP, err := godaddyndns.GetExternalIP()

	if err != nil {
		return err
	}

	dnsIP, err := godaddyndns.GetGoDaddyGetAddress(cmdConfig.Key, cmdConfig.Secret, cmdConfig.Domain, cmdConfig.Host)

	if err != nil {
		return err
	}

	if dnsIP.Equal(extIP) {
		return nil
	}

	log.Printf("updating %s.%s address from %s to %s", cmdConfig.Host, cmdConfig.Domain, dnsIP.String(), extIP.String())

	err = godaddyndns.GetGoDaddySetAddress(cmdConfig.Key, cmdConfig.Secret, cmdConfig.Domain, cmdConfig.Host, extIP)

	if err != nil {
		return err
	}

	log.Printf("address of %s.%s has been successfully updated to %s", cmdConfig.Host, cmdConfig.Domain, extIP.String())

	return nil
}

func main() {
	log.SetFlags(0)

	flag.Parse()

	if cmdFlags.ConfigFilePath == "" {
		log.Fatalln("configuration file was not specified")
	}

	var err error

	cmdConfig, err = cmd.LoadConfiguration(cmdFlags.ConfigFilePath)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("starting to monitor external address for %s.%s using %s", cmdConfig.Host, cmdConfig.Domain, godaddyndns.GetExternalIPService())

	for {
		go func() {
			if err := CheckAndUpdate(); err != nil {
				log.Println(err)
			}
		}()

		time.Sleep(time.Duration(cmdFlags.CheckInterval) * time.Second)
	}
}
