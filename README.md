# dyndns-cloudflare-adapter
HTTP server implementing the DynDNS.org [protocol](https://help.dyn.com/remote-access-api/perform-update/)
allowing to update the CloudFront DNS records with dynamic IP directly from your router.

## Building and running
1. `docker build -t dyndns-cloudflare-adapter .`
2. `docker run -it -p 8080:8080 -e CF_API_EMAIL=YourEmail -e CF_API_KEY=YourGlobalAPIKey dyndns-cloudflare-adapter`
3. Configure your router to call the server using custom endpoint.
   It needs to provide `hostname` and `myip` query string parameters:
    - `myip` is the new IP address of the server
    - `hostname` is a regexp that is used to select domain names from your account. 
      Only domains matching the regexp will be updated with the provided IP.
