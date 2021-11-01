# DnsMon

###Quick and dirty DNS monitor for NewRelic


##Required Environment Variables

NEWRELIC_KEY=Your NewRelic Key

MONITOR_NAME=NICKT1 Monitor Name

POLL_INTERVAL=45 Number of seconds between polls

dns.host.0=www.newrelic.com

dns.host.1=www.google.com

...

dns.host.99= Up to 100 hosts to monitor

##Run in Docker
Copy .env.template to .env Add domains to monitor, a monitor name and a New Relic license key

docker build . -t dnsmon

docker run --env-file .env dnsmon

##Generated events
####dns_mon:

dns_error set to 1 if error

duration set to ms for DNS lookup

####DNSMonTick:
Metadata about state of DNS monitor sent evey three minutes
