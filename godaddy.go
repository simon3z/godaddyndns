package godaddyndns

// cspell:ignore godaddyndns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type HostResponse struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

func GoDaddySetAddress(key, secret, domain, host string, address net.IP) error {
	u := []HostResponse{{Data: address.String()}}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(u)

	req, err := http.NewRequest("PUT", GoDaddyDomainAPIEndpoint(domain, host), body)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", GoDaddyAuthorization(key, secret))
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	io.Copy(os.Stdout, res.Body)

	return nil
}

func GoDaddyDomainAPIEndpoint(domain, host string) string {
	return fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, host)
}

func GoDaddyAuthorization(key, secret string) string {
	return fmt.Sprintf("sso-key %s:%s", key, secret)
}
