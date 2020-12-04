package main

import (
	"./cf"
	"regexp"
)

type HostnamePatternError struct {
	cause error
}
func (err HostnamePatternError) Error() string {
	return "hostname pattern invalid"
}
func (err HostnamePatternError) Unwrap() error {
	return err.cause
}

type DNSProviderError struct {
	cause error
	message string
}
func (err DNSProviderError) Error() string {
	return err.message
}
func (err DNSProviderError) Unwrap() error {
	return err.cause
}

type Updater struct {
	cf *cf.Client
}

func (u Updater) UpdateDNS(newIP string, hostnamePattern string) error {
	var err error
	var records []cf.DNSRecord
	var pattern *regexp.Regexp

	if hostnamePattern != "" {
		pattern, err = regexp.Compile(hostnamePattern)
		if err != nil {
			return HostnamePatternError{err}
		}
	}

	records, err = u.cf.FindDNSRecords(pattern)
	if err != nil {
		return DNSProviderError{err, "cannot find DNS records"}
	}

	for _, record := range records {
		if record.IP != newIP {
			if err = u.cf.UpdateDNSRecord(record, newIP); err != nil {
				return DNSProviderError{err, "failed updating DNS record"}
			}
		}
	}

	return nil
}
