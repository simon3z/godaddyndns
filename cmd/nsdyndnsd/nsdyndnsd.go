package main

// cspell:ignore nsdyndns

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/simon3z/nsdyndns"
	"github.com/simon3z/nsdyndns/cmd"
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
	extIP, err := nsdyndns.GetExternalIP()

	if err != nil {
		return err
	}

	dnsAddrs, err := net.LookupIP(cmdConfig.FullDomain())

	if err != nil {
		return fmt.Errorf("address lookup for %s.%s failed: %w", cmdConfig.Host, cmdConfig.Domain, err)
	}

	if len(dnsAddrs) != 1 {
		return fmt.Errorf("unexpected multiple ip addresses found: %#v", dnsAddrs)
	}

	if dnsAddrs[0].Equal(extIP) {
		return nil
	}

	log.Printf("new external ip address detected: %s", extIP.String())

	log.Printf("updating %s.%s address from %s to %s", cmdConfig.Host, cmdConfig.Domain, dnsAddrs[0].String(), extIP.String())

	err = nsdyndns.GoDaddySetAddress(cmdConfig.Key, cmdConfig.Secret, cmdConfig.Domain, cmdConfig.Host, extIP)

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

	log.Printf("starting to monitor external address for %s.%s using %s", cmdConfig.Host, cmdConfig.Domain, nsdyndns.GetExternalIPService())

	for {
		go func() {
			if err := CheckAndUpdate(); err != nil {
				log.Println(err)
			}
		}()

		time.Sleep(time.Duration(cmdFlags.CheckInterval) * time.Second)
	}
}
