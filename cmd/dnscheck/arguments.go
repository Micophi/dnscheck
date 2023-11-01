package main

import "fmt"

var version = "devel"

type CliArgs struct {
	Domains    string `arg:"--domains" help:"Path to file with list of domains to check"`
	DnsConfigs string `arg:"--config" help:"Override the config file used with the one provided"`
	Output     string `arg:"--output" help:"Override default output path"`
}

func (CliArgs) Version() string {
	return fmt.Sprintf("DNSCheck version %s\n", version)
}

func (CliArgs) Description() string {
	return "\nDNSCheck is a tool to test a list of domains against multiple DNS server simultaneously and return the number that was blocked.\n"
}

func (CliArgs) Epilogue() string {
	return "For more information visit: https://github.com/Micophi/dnscheck\n"
}
