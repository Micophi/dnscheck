package structs

import (
	"dnscheck/internal/utilities"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/goccy/go-yaml"

	"go.uber.org/atomic"
)

type DnsServer struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Ip          string            `yaml:"ip"`
	RateLimit   int               `yaml:"rateLimit"`
	Count       atomic.Int32      `default:"0"`
	Blocked     atomic.Int32      `default:"0"`
	Retries     atomic.Int32      `default:"0"`
	Skip        atomic.Int32      `default:"0"`
	AvgRtt      atomic.Duration   `default:"0"`
	Client      upstream.Upstream `yaml:"-" json:"-"`
}

func (d *DnsServer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name        string
		Description string
		Ip          string
		RateLimit   int
		Count       int
		Blocked     int
		Retries     int
		Skip        int
		AvgRtt      string
	}{
		Name:        d.Name,
		Description: d.Description,
		Ip:          d.Ip,
		RateLimit:   d.RateLimit,
		Count:       int(d.Count.Load()),
		Blocked:     int(d.Blocked.Load()),
		Retries:     int(d.Retries.Load()),
		Skip:        int(d.Skip.Load()),
		AvgRtt:      fmt.Sprint(d.AvgRtt.Load()),
	})
}

type DnsServers struct {
	DnsServers []DnsServer `yaml:"dnsservers"`
	RateLimit  int         `yaml:"rateLimit"`
}

func (dnsServers DnsServers) Save(outputpath string) {
	results, err := json.Marshal(&dnsServers)
	utilities.CheckError(err)
	results, err = yaml.JSONToYAML(results)
	utilities.CheckError(err)

	err = os.WriteFile(outputpath, results, 0644)
	utilities.CheckError(err)
}

func (dnsServers DnsServers) SaveDefault() {
	filename := fmt.Sprintf("%s.yaml", time.Now().Format("2006-01-02_15h04m05"))
	dnsServers.Save(filename)
}

func (dnsServers DnsServers) PrintSummary() {
	fmt.Println("############################################ SUMMARY ###########################################")
	for _, dnsServer := range dnsServers.DnsServers {
		fmt.Println("Blocked:", dnsServer.Blocked.Load(), "\t|", "Total:", dnsServer.Count.Load(), "\t|", "Skipped: ", dnsServer.Skip.Load(), "\t|", "AvgRtt:", dnsServer.AvgRtt.Load(), "\t|", "Name:", dnsServer.Name)
	}
	fmt.Println("################################################################################################")

}
