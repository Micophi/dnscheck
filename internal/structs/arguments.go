package structs

type CliArgs struct {
	Domains    string `arg:"--domains"`
	DnsConfigs string `arg:"--config"`
	Output     string `arg:"--output"`
}
