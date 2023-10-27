package main

import (
	"context"
	"dnscheck/internal/structs"
	"dnscheck/internal/utilities"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/time/rate"
)

var GLOBAL_RATE_LIMIT = 25
var GLOBAL_MAX_RETRIES = 20

func Lookup(domain string, dnsServer string) (bool, time.Duration, int, bool) {
	msg := new(dns.Msg)
	msg.SetQuestion(fmt.Sprintf("%s.", domain), dns.TypeA)

	client := dns.Client{}
	retries := 0
	var err error
	for err == nil {
		in, rtt, err := client.Exchange(msg, fmt.Sprintf("%s:53", dnsServer))

		if err != nil {
			// os.Stderr.WriteString(fmt.Sprintf("|%s| Error while resolving: %s. Retrying....\n", dnsServer, domain))
			retries++
		} else if in.Answer == nil || len(in.Answer) == 0 {
			return true, rtt, retries, false
		} else {
			var ip string
			for _, element := range in.Answer {
				if _, ok := element.(*dns.A); ok {
					ip = element.(*dns.A).A.String()
					break
				}
			}
			if ip == "0.0.0.0" {
				return true, rtt, retries, false
			}

			return false, rtt, retries, false
		}
		if retries == GLOBAL_MAX_RETRIES {
			return false, 0, retries, true
		}
	}
	return false, 0, retries, false
}

func computeAvgRtt(latencies []time.Duration, dnsServer structs.DnsServer) time.Duration {
	// calculate average latency
	total := time.Duration(0)
	for _, latency := range latencies {
		total += latency
	}
	total = total / (time.Duration(dnsServer.Count - dnsServer.Skip))
	return total
}

func testDNS(domains []string, dnsServer *structs.DnsServer, wg *sync.WaitGroup, progressBar *mpb.Bar) {
	if dnsServer.RateLimit == 0 {
		dnsServer.RateLimit = GLOBAL_RATE_LIMIT
	}
	limiter := rate.NewLimiter(rate.Limit(dnsServer.RateLimit), 0)
	defer wg.Done()
	latencies := []time.Duration{}
	for _, domain := range domains {
		limiter.Wait(context.Background())
		sinkholed, rtt, retries, timedout := Lookup(domain, dnsServer.Ip)
		dnsServer.Retries += retries

		if timedout {
			dnsServer.Skip++
		} else {
			latencies = append(latencies, rtt)
			dnsServer.Count++
			if sinkholed {
				dnsServer.Blocked++
			}
		}
		progressBar.Increment()
	}
	dnsServer.AvgRtt = computeAvgRtt(latencies, *dnsServer)
}

func readConfigurations() {
	viper.SetConfigName("dnscheck")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/dnscheck/")
	err := viper.ReadInConfig() // Find and read the config file
	utilities.CheckError(err)
}

func getDomainsFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	utilities.CheckError(err)

	lines := strings.Split(string(content), "\n")
	return lines
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
		readConfigurations()
	}

	domains := getDomainsFromFile(args.Domains)
	start := time.Now()

	dnsServers := structs.DnsServers{}
	err := viper.Unmarshal(&dnsServers)
	utilities.CheckError(err)

	var wg sync.WaitGroup
	progressWaitGroup := mpb.New(mpb.WithWaitGroup(&wg))

	total := len(domains)
	for index, dnsServer := range dnsServers.DnsServers {

		name := fmt.Sprintf("%s:", dnsServer.Name)
		bar := progressWaitGroup.AddBar(int64(total),
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
		wg.Add(1)
		go testDNS(domains, &dnsServers.DnsServers[index], &wg, bar)

	}
	progressWaitGroup.Wait()

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
