package structs

import (
	"dnscheck/internal/utilities"
	"fmt"
	"os"
	"time"

	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
)

type DnsServer struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Ip          string          `yaml:"ip"`
	RateLimit   int             `yaml:"rateLimit"`
	Count       atomic.Int32    `default:"0"`
	Blocked     atomic.Int32    `default:"0"`
	Retries     atomic.Int32    `default:"0"`
	Skip        atomic.Int32    `default:"0"`
	AvgRtt      atomic.Duration `default:"0"`
}

type DnsServers struct {
	DnsServers []DnsServer `yaml:"dnsservers"`
	RateLimit  int         `yaml:"rateLimit"`
}

func (dnsServers DnsServers) Save(outputpath string) {
	results, err := yaml.Marshal(&dnsServers)
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
