package cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const baseURL = "https://api.cloudflare.com/client/v4/"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type Client struct {
	header http.Header
	zones []string
}

type DNSRecord struct {
	ID string `json:"id"`
	Name string `json:"name"`
	IP string `json:"content"`
	Zone string `json:"zone_id"`
}

func NewClient(email, key string) (*Client, error) {
	client := &Client{
		header: http.Header{
			"Content-Type": {"application/json"},
			"X-Auth-Email": {email},
			"X-Auth-Key": {key},
		},
	}

	r, err := client.get("zones?status=active&per_page=50")
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}

	if result, ok := data["result"]; ok {
		result := result.([]interface {})
		zones := make([]string, 0, len(result))

		for _, zone := range result {
			zone := zone.(map[string]interface{})
			zones = append(zones, zone["id"].(string))
		}

		client.zones = zones
	}

	if client.zones == nil || len(client.zones) == 0 {
		return nil, fmt.Errorf("no zones found")
	}

	return client, nil
}

func (cf *Client) FindDNSRecords(pattern *regexp.Regexp) ([]DNSRecord, error) {
	allRecords := make([]DNSRecord, 0)

	for _, zone := range cf.zones {
		records, err := cf.fetchRecords(zone, pattern)
		if err != nil {
			return allRecords, err
		}
		allRecords = append(allRecords, records...)
	}

	return allRecords, nil
}

func (cf *Client) UpdateDNSRecord(record DNSRecord, ip string) error {
	body, _ := json.Marshal(map[string]string{
		"type": "A",
		"name": record.Name,
		"content": ip,
	})

	r, err := cf.put(fmt.Sprintf("zones/%s/dns_records/%s", record.Zone, record.ID), body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode / 100 != 2 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Updating DNS record failed:\n%w", fmt.Errorf(string(body)))
	}

	return nil
}

func (cf *Client) put(path string, body []byte) (*http.Response, error) {
	r, err := http.NewRequest("PUT", baseURL + path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	r.Header = cf.header
	return httpClient.Do(r)
}

func (cf *Client) get(path string) (*http.Response, error) {
	r, err := http.NewRequest("GET", baseURL + path, nil)
	if err != nil {
		return nil, err
	}
	r.Header = cf.header
	return httpClient.Do(r)
}

func decodeRecords(body io.Reader) ([]DNSRecord, error)  {
	decoder := json.NewDecoder(body)

	var err error
	var token json.Token

	// skip until we get to "result" array
	for decoder.More() && token != "result" {
		token, err = decoder.Token()
		if err != nil {
			return nil, err
		}
	}

	// decode each object inside array into slice of DNSRecord struct
	results := make([]DNSRecord, 0)
	if err = decoder.Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}

func (cf *Client) fetchRecords(zone string, pattern *regexp.Regexp) ([]DNSRecord, error) {
	r, err := cf.get(fmt.Sprintf("zones/%s/dns_records?type=A&per_page=100", zone))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode / 100 != 2 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Fetching DNS records failed:\n%w", fmt.Errorf(string(body)))
	}

	records, err := decodeRecords(r.Body)
	if err != nil {
		return nil, err
	}

	if pattern == nil {
		return records, nil
	}

	matchedRecords := make([]DNSRecord, 0, len(records))
	for _, record := range records {
		if pattern.MatchString(record.Name) {
			matchedRecords = append(matchedRecords, record)
		}
	}

	return matchedRecords, nil
}