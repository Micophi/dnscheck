package main

import (
	"dnscheck/internal/dnscheck"
	"dnscheck/internal/structs"
	"dnscheck/internal/utilities"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/alexflint/go-arg"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func createProgressBars(dnsServers []structs.DnsServer, length int, progressWaitGroup *mpb.Progress) map[int]*mpb.Bar {
	var progressBars = make(map[int]*mpb.Bar)

	for index, dnsServer := range dnsServers {
		name := fmt.Sprintf("%s:", dnsServer.Name)
		progressBars[index] = progressWaitGroup.AddBar(int64(length),
			mpb.PrependDecorators(
				decor.Name(name),
				decor.CountersNoUnit(" %d/%d", decor.WCSyncWidth),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace), "done",
				),
			),
		)
	}
	return progressBars
}

var version = "devel"

type CliArgs struct {
	Domains    string   `arg:"positional" help:"path to file with list of domains to check"`
	ExtraDNS   []string `arg:"--dns,separate" help:"override config file and test de provided DNS servers. Can be used multiple times."`
	DnsConfigs string   `arg:"--config" help:"override the config file used with the one provided"`
	Output     string   `arg:"--output" help:"override default output path"`
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

func main() {

	var args CliArgs
	p := arg.MustParse(&args)

	if args.Domains == "" {
		p.Fail("You must provide --domains")
	}

	if args.DnsConfigs != "" {
		viper.SetConfigType("yaml")
		file, err := os.Open(args.DnsConfigs)
		utilities.CheckError(err)
		viper.ReadConfig(file)

	} else {
		utilities.ReadConfigurations()
	}

	domains := utilities.GetDomainsFromFile(args.Domains)
	start := time.Now()

	dnsServers := structs.DnsServers{}

	if len(args.ExtraDNS) > 0 {
		for _, dns := range args.ExtraDNS {
			cliDNS := structs.DnsServer{Name: dns, Ip: dns, RateLimit: 25}
			dnsServers.DnsServers = append(dnsServers.DnsServers, cliDNS)
		}
	} else {
		err := viper.Unmarshal(&dnsServers)
		utilities.CheckError(err)
	}

	var wg sync.WaitGroup
	progressWaitGroup := mpb.New(mpb.WithWaitGroup(&wg))

	var progressBars = createProgressBars(dnsServers.DnsServers[:], len(domains), progressWaitGroup)
	var rateLimiters = dnscheck.CreateRateLimiters(dnsServers.DnsServers[:], dnsServers.RateLimit)

	for index := range dnsServers.DnsServers {
		var err error
		dnsServers.DnsServers[index].Client, err = upstream.AddressToUpstream(dnsServers.DnsServers[index].Ip, &upstream.Options{})
		utilities.CheckError(err)
	}

	for _, domain := range domains {
		for index := range dnsServers.DnsServers {
			wg.Add(1)
			go dnscheck.DomainNameCheck(domain, &dnsServers.DnsServers[index], &wg, progressBars[index], rateLimiters[index])
		}
	}

	progressWaitGroup.Wait()
	dnscheck.UpdateAverageRtt(dnsServers.DnsServers[:])

	// End of run summary
	fmt.Println("")
	fmt.Println("Completed in: ", time.Since(start))
	fmt.Println("")
	if args.Output == "" {
		dnsServers.SaveDefault()
	} else {
		dnsServers.Save(args.Output)
	}

	dnsServers.PrintSummary()
}
