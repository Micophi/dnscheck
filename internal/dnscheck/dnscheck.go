package dnscheck

import (
	"context"
	"dnscheck/internal/structs"
	"fmt"
	"sync"
	"time"

	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
	"github.com/vbauerster/mpb/v8"
	"golang.org/x/time/rate"
)

var sinkholeIp = "0.0.0.0"
var adguardRedirectIp = "94.140.14.33"
var ciraRedirectIps = []string{"75.2.78.236", "99.83.179.4", "99.83.178.7", "75.2.110.227"}
var blockedIpAnswers = append(ciraRedirectIps, sinkholeIp, adguardRedirectIp)

var maximumRetries = int32(20)

func IsDomainBlocked(domain string, client upstream.Upstream) (bool, time.Duration, int32, bool) {
	msg := new(dns.Msg)
	msg.SetQuestion(fmt.Sprintf("%s.", domain), dns.TypeA)

	retries := int32(0)
	var err error
	for err == nil {

		start := time.Now()
		in, err := client.Exchange(msg)
		var rtt = time.Since(start)
		if err != nil {
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
			for _, blocked := range blockedIpAnswers {
				if ip == blocked {
					return true, rtt, retries, false
				}
			}

			return false, rtt, retries, false
		}
		if retries == maximumRetries {
			return false, 0, retries, true
		}
	}
	return false, 0, retries, false
}

func DomainNameCheck(domain string, dnsServer *structs.DnsServer, wg *sync.WaitGroup, progressBar *mpb.Bar, rateLimiter *rate.Limiter) {
	rateLimiter.Wait(context.Background())

	sinkholed, rtt, retries, timeout := IsDomainBlocked(domain, dnsServer.Client)
	dnsServer.Retries.Add(retries)

	if timeout {
		dnsServer.Skip.Inc()
	} else {
		dnsServer.AvgRtt.Add(rtt)
		dnsServer.Count.Inc()
		if sinkholed {
			dnsServer.Blocked.Inc()
		}
	}
	progressBar.Increment()
	wg.Done()

}

func CreateRateLimiters(dnsServers []structs.DnsServer, globalRateLimit int) map[int]*rate.Limiter {
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

func UpdateAverageRtt(dnsServers []structs.DnsServer) {
	for index, dnsServer := range dnsServers {
		if dnsServer.Count.Load() > 0 {
			dnsServers[index].AvgRtt.Swap(dnsServer.AvgRtt.Load() / time.Duration(dnsServer.Count.Load()))
		}
	}
}
