package nsdyndns

// cspell:ignore nsdyndns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type GoDaddyService struct {
	Key    string
	Secret string
}

type HostResponse struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

func NewGoDaddyService(key, secret string) NameService {
	return &GoDaddyService{key, secret}
}

func (s *GoDaddyService) SetAddress(domain, host string, address net.IP) error {
	u := []HostResponse{{Data: address.String()}}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(u)

	req, err := http.NewRequest("PUT", goDaddyDomainAPIEndpoint(domain, host), body)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", goDaddyAuthorization(s.Key, s.Secret))
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("couldn't update ip address: %s", strings.TrimSpace(string(resBody)))
	}

	return nil
}

func goDaddyDomainAPIEndpoint(domain, host string) string {
	return fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, host)
}

func goDaddyAuthorization(key, secret string) string {
	return fmt.Sprintf("sso-key %s:%s", key, secret)
}
