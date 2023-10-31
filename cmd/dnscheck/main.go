package main

import (
	"context"
	"dnscheck/internal/dnscheck"
	"dnscheck/internal/structs"
	"dnscheck/internal/utilities"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/time/rate"
)

var GLOBAL_RATE_LIMIT = 25

func domainNameCheck(domain string, dnsServer *structs.DnsServer, wg *sync.WaitGroup, progressBar *mpb.Bar, rateLimiter *rate.Limiter) {
	rateLimiter.Wait(context.Background())

	sinkholed, rtt, retries, timedout := dnscheck.IsDomainBlocked(domain, dnsServer.Ip)
	dnsServer.Retries += retries

	if timedout {
		dnsServer.Skip++
	} else {
		dnsServer.AvgRtt += rtt
		dnsServer.Count += 1
		if sinkholed {
			dnsServer.Blocked++
		}
	}

	wg.Done()
	progressBar.Increment()
}

func updateAverageRtt(dnsServers []structs.DnsServer) {
	for index, dnsServer := range dnsServers {
		dnsServers[index].AvgRtt = dnsServer.AvgRtt / time.Duration(dnsServer.Count)
	}
}

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

func createRateLimiters(dnsServers []structs.DnsServer, globalRateLimit int) map[int]*rate.Limiter {
	var rateLimiters = make(map[int]*rate.Limiter)
	var globalLimiter = rate.NewLimiter(rate.Limit(globalRateLimit), 1)

	for index, dnsServer := range dnsServers {
		if dnsServer.RateLimit > 0 {
			rateLimiters[index] = rate.NewLimiter(rate.Limit(dnsServer.RateLimit), 1)
		} else if dnsServer.RateLimit == 0 {
			rateLimiters[index] = globalLimiter
		}
	}
	return rateLimiters
}

func main() {

	var args struct{ structs.CliArgs }
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
	err := viper.Unmarshal(&dnsServers)
	utilities.CheckError(err)

	var wg sync.WaitGroup
	progressWaitGroup := mpb.New(mpb.WithWaitGroup(&wg))

	var progressBars = createProgressBars(dnsServers.DnsServers[:], len(domains), progressWaitGroup)
	var rateLimiters = createRateLimiters(dnsServers.DnsServers[:], dnsServers.RateLimit)

	for _, domain := range domains {
		for index := range dnsServers.DnsServers {
			wg.Add(1)
			go domainNameCheck(domain, &dnsServers.DnsServers[index], &wg, progressBars[index], rateLimiters[index])
		}
	}

	progressWaitGroup.Wait()
	updateAverageRtt(dnsServers.DnsServers[:])

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
