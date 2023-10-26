# Dnscheck

This project was inspired by [this post](https://forums.lawrencesystems.com/t/which-is-the-best-dns-for-secure-browsing-cloudflare-quad9-nextdns-and-adguard-dns-youtube-release/18910) from Lawrence of Lawrence Systems.

![](examples/demo.gif)

## Configuration file
Dnscheck will lock for a `dnscheck.yaml` configuration file in the same directory has it is being run. This file must contain the DNS servers you want to test. There is a sample config file available under `/configs` with preconfigured DNS servers. 

Currently, dnscheck only supports DNS over UDP on port 53.

```yaml
dnsservers:
  - name: Quad9   # Descriptive name of the DNS server
    description: Malware Blocking, DNSSEC Validation
    ip: 9.9.9.9   # The IP of the DNS server
    rateLimit: 15 # Limit the number of queries per second to send to the DNS server
```

## Behavior
Currently the application retries to resolve a domain up to 20 times before skipping it. In my limited testing, the queries sometime fails probably due to my test configurations.

## Results
By default, `dnscheck` will output the results in the current folder using the current date and time as filename(i.e. YYYY-MM-DD_HHhMMmSS.yaml). The output has the same format as the config file with only more fields added to it for details.

Example:
```yaml
dnsservers:
    - name: Quad9
      description: Malware Blocking, DNSSEC Validation
      ip: 9.9.9.9
      rateLimit: 15
      count: 101          # Total number of domain queried
      blocked: 101        # Total number of domain blocked
      retries: 0          # Number of retries executed accross all domains queried
      skip: 0             # Number of domain skipped due to too many retries
      avgrtt: 42.312169ms # Average round trip time of the queries
```

## Scripts
The script folder has some simple script to download domain list to test. They require `wget` and `sponge` to run properly.

```bash
sudo apt-get install -y wget moreutils
```


