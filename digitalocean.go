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

type DigitalOceanService struct {
	Token string
}

type DomainRecord struct {
	Id   int64
	Type string
	Name string
	Data string
	Ttl  int
}

const (
	DigitalOceanDefaultTtl = 60
)

func NewDigitalOceanService(token string) NameService {
	return &DigitalOceanService{token}
}

func NewDigitalOceanRequest(method, path, token string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://api.digitalocean.com/v2/%s", path), body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func DigitalOceanGetRecordId(token, domain, host string) (int64, error) {
	req, err := NewDigitalOceanRequest("GET", fmt.Sprintf("domains/%s/records", domain), token, nil)

	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("couldn't update ip address: %s", strings.TrimSpace(string(resBody)))
	}

	address := struct {
		DomainRecords []DomainRecord `json:"domain_records"`
	}{}

	err = json.Unmarshal(resBody, &address)

	if err != nil {
		return 0, err
	}

	for i := range address.DomainRecords {
		if address.DomainRecords[i].Name == host {
			return address.DomainRecords[i].Id, nil
		}
	}

	return 0, nil
}

func (s *DigitalOceanService) SetAddress(domain, host string, address net.IP) error {
	id, err := DigitalOceanGetRecordId(s.Token, domain, host)

	if err != nil {
		return err
	}

	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(DomainRecord{Data: address.String(), Ttl: DigitalOceanDefaultTtl})

	req, err := NewDigitalOceanRequest("PUT", fmt.Sprintf("domains/%s/records/%d", domain, id), s.Token, reqBody)

	if err != nil {
		return err
	}

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

	resAddress := struct {
		DomainRecords DomainRecord `json:"domain_record"`
	}{}

	err = json.Unmarshal(resBody, &resAddress)

	if err != nil {
		return err
	}

	if resAddress.DomainRecords.Data != address.String() || resAddress.DomainRecords.Ttl != DigitalOceanDefaultTtl {
		return fmt.Errorf("possible failure in domain update: %#v", resAddress)
	}

	return nil
}
