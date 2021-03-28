package godaddyndns

// cspell:ignore godaddyndns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type HostResponse struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

func GetGoDaddyGetAddress(key, secret, domain, host string) (net.IP, error) {
	req, err := http.NewRequest("GET", GoDaddyDomainAPIEndpoint(domain, host), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", GoDaddyAuthorization(key, secret))
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("couldn't retrieve ip address: %s", body)
	}

	address := []HostResponse{}

	err = json.Unmarshal(body, &address)

	if err != nil {
		return nil, err
	}

	if len(address) == 0 {
		return nil, fmt.Errorf("couldn't find any entry for %s.%s", host, domain)
	}

	for _, v := range address {
		if v.Name != host {
			continue
		}

		ip := net.ParseIP(v.Data)

		if ip != nil {
			return ip, nil
		}
	}

	return nil, fmt.Errorf("none of the entries matched %s.%s", host, domain)
}

func GetGoDaddySetAddress(key, secret, domain, host string, address net.IP) error {
	u := []HostResponse{{Data: address.String()}}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(u)

	req, err := http.NewRequest("PUT", GoDaddyDomainAPIEndpoint(domain, host), body)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", GoDaddyAuthorization(key, secret))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	io.Copy(os.Stdout, resp.Body)

	return nil
}

func GoDaddyDomainAPIEndpoint(domain, host string) string {
	return fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, host)
}

func GoDaddyAuthorization(key, secret string) string {
	return fmt.Sprintf("sso-key %s:%s", key, secret)
}
