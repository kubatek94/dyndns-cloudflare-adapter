# dyndns-cloudflare-adapter
HTTP server implementing the DynDNS.org [protocol](https://help.dyn.com/remote-access-api/perform-update/),
or a STUN daemon allowing to update the CloudFlare DNS records with a dynamic IP.

## Building
`docker build -t dyndns-cloudflare-adapter .`


---


## Running as daemon (with STUN IP detection)

`docker run -it -e CF_API_EMAIL=YourEmail -e CF_API_KEY=YourGlobalAPIKey -e HOSTNAME=.+\-domain\.co\.uk dyndns-cloudflare-adapter`


---

## Running as DynDNS server
Configure your router to call the server using custom endpoint.
It needs to provide `hostname` and `myip` query string parameters:
 - `myip` is the new IP address of the server
 - `hostname` is a regexp that is used to select domain names from your account.
   Only domains matching the regexp will be updated with the provided IP.

`docker run -it -p 8080:8080 -e CF_API_EMAIL=YourEmail -e CF_API_KEY=YourGlobalAPIKey dyndns-cloudflare-adapter`


### Example OpenWrt config
In `/etc/config/ddns`
```text
config ddns 'global'
    option ddns_dateformat '%F %R'
    option ddns_loglines '250'
    option upd_privateip '0'
    option use_curl '1'

config service 'myddns_ipv4'
    option enabled '1'
    option interface 'wan'
    option ip_source 'network'
    option ip_network 'wan'
    option use_https '1'
    option cacert 'IGNORE'
    option force_interval '1'
    option force_unit 'days'
    option lookup_host 'your-domain.com' # used to check if the current IP is already up to date (nslookup for DNS)
    option domain '.+-domain\.com' # the regexp of domains to update
    option username 'HttpBasicUser'
    option password 'HttpBasicPass'
    option dns_server '1.1.1.1' # optional
    option update_url 'http://[USERNAME]:[PASSWORD]@dyndns.your-domain.com/nic/update?hostname=[DOMAIN]&myip=[IP]'
```
