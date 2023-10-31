package structs

import (
	"dnscheck/internal/utilities"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type DnsServer struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Ip          string        `yaml:"ip"`
	RateLimit   int           `yaml:"rateLimit"`
	Count       int           `default:"0"`
	Blocked     int           `default:"0"`
	Retries     int           `default:"0"`
	Skip        int           `default:"0"`
	AvgRtt      time.Duration `default:"0"`
}

type DnsServers struct {
	DnsServers []DnsServer `yaml:"dnsservers"`
	RateLimit  int         `yaml:"rateLimit"`
}

func (dnsServers DnsServers) Save(outputpath string) {
	results, err := yaml.Marshal(&dnsServers)
	utilities.CheckError(err)
	// filename := fmt.Sprintf("%s.yaml", time.Now().Format("2006-01-02_15h04m05"))

	err = os.WriteFile(outputpath, results, 0644)
	utilities.CheckError(err)
}

func (dnsServers DnsServers) SaveDefault() {
	filename := fmt.Sprintf("%s.yaml", time.Now().Format("2006-01-02_15h04m05"))
	dnsServers.Save(filename)
	// results, err := yaml.Marshal(&dnsServers)
	// utilities.CheckError(err)
	// filename := fmt.Sprintf("%s.yaml", time.Now().Format("2006-01-02_15h04m05"))

	// err = os.WriteFile(filename, results, 0644)
	// utilities.CheckError(err)
}

func (dnsServers DnsServers) PrintSummary() {
	fmt.Println("############################################ SUMMARY ###########################################")
	for _, dnsServer := range dnsServers.DnsServers {
		fmt.Println("Blocked:", dnsServer.Blocked, "\t|", "Total:", dnsServer.Count, "\t|", "Skipped: ", dnsServer.Skip, "\t|", "AvgRtt:", dnsServer.AvgRtt, "\t|", "Name:", dnsServer.Name)
	}
	fmt.Println("################################################################################################")

}
