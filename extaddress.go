package godaddyndns

// cspell:ignore godaddyndns ipify

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var externalIPServiceURL = url.URL{Scheme: "https", Host: "api.ipify.org"}

func GetExternalIP() (net.IP, error) {
	resp, err := http.Get(externalIPServiceURL.String())

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("couldn't retrieve ip address: %s", body)
	}

	ip := net.ParseIP(strings.TrimSpace(string(body)))

	if ip == nil {
		return nil, fmt.Errorf("couldn't parse ip address %s", body)
	}

	return ip, nil
}

func GetExternalIPService() string {
	return externalIPServiceURL.Host
}
