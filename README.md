# Dnscheck

This project was inspired by [this post](https://forums.lawrencesystems.com/t/which-is-the-best-dns-for-secure-browsing-cloudflare-quad9-nextdns-and-adguard-dns-youtube-release/18910) from Lawrence of Lawrence Systems.

Dnscheck is a tool to check if a DNS server will properly resolve or block a provided domain name. Some DNS servers will allow most name resolution except widely known malicious domains, whereas others will block various categories (e.g. Ads, trackers, adult content, compromised domains, etc).

![](examples/demo.gif)

## Installation

The release page contains the binary and provided configuration file. 

### Debian-based OS (Debian, Ubuntu, etc.)
For debian-based operating systems, a .deb package will install the binary and deploy the configuration file under `/etc/dnscheck`
Also the scripts will be deployed on the system under `/usr/bin`.


### Installation via packagecloud

Debian packages are [hosted on packagecloud](https://packagecloud.io/Micophi/dnscheck). There is a user script to install the repository:
```bash
curl -s https://packagecloud.io/install/repositories/Micophi/dnscheck/script.deb.sh | sudo bash
```

If you would like to validate the script and manually add the repository, [there is a manual installation documentation.](https://packagecloud.io/Micophi/dnscheck/install#manual-deb)

## Usage

```bash
Usage: dnscheck [--dns DNS] [--config CONFIG] [--output OUTPUT] [DOMAINS]

Positional arguments:
  DOMAINS                path to file with list of domains to check

Options:
  --dns DNS              override config file and test de provided DNS servers. Can be used multiple times.
  --config CONFIG        override the config file used with the one provided
  --output OUTPUT        override default output path
  --help, -h             display this help and exit
  --version              display version and exit
```

## Configuration file
Dnscheck will look for a `dnscheck.yaml` configuration file in the same directory has it is being run. This file must contain the DNS servers you want to test. There is a sample config file available under `/configs` with preconfigured DNS servers. 

`dnscheck` supports DoT and DoH as well as traditional DNS query.


```yaml
rateLimit: 50 # Global rate limit (Shared between all servers tested)
dnsservers:
  - name: Quad9 DoH  # Descriptive name of the DNS server
    description: Malware Blocking, DNSSEC Validation
    ip: https://dns.quad9.net/dns-query   # The IP or address of the DNS server
    rateLimit: 15 # (Overrides the global rate limit for this server)Limit the number of queries per second to send to the DNS server
  - name: Cloudflare
    description: Cloudflare DNS(https://1.1.1.1/)
    ip: 1.1.1.1
```

## Behavior
Currently the application retries to resolve a domain up to 20 times before skipping it. In my limited testing, the queries sometime fails due to network configurations (e.g. IDS, internet providers, services rate limit, etc.)

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
The script folder has some utility scripts to download domain lists to test. They require `wget` and `sponge` to run properly.

```bash
sudo apt-get install -y wget moreutils
```



