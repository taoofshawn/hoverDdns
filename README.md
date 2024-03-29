# hoverDdns - Broken
## I believe Hover disabled this API to make room for a new [partners product](https://partners.hover.com/reference/about-hover)
## Hopefully its something else but can't connect for now.

Keep a hover DNS record up to date with your external IP. 
A use-case would be a home LAN with the following features:
- NAT router - with a public WAN ip that may change periodically
- Some NAS device or other "always on" computer capable of running Docker
- Some desire to always know the external IP (i.e. VPN)
- you already have a domain on Hover for some other purpose

## About
I have been using one of the free dynamic dns providers with a monthly nag email and didnt hear about duckdns until i got too far and interested in finishing this one :/

## Why another hover updater
This is based on [hover-dns-updater.](https://github.com/texasaggie97/hover-dns-updater). I tried to use that as a docker container and it didn't work for me right away on the first try so I thought I'd figure it out by making my own!  
I removed everything that wasn't needed for a docker implementation  
Everything was working until I spun up the first container and I ended up with the same error that the original one had!  It turned out to be some SSL issue between the newer python base image and hover.  Adding an ssl package to the requirements.txt file ended up solving it.  
I will keep this as a learning project and continually optimize and improve.


## Setup
This requires at least one registered hover domain.  
If you haven't already, login and create an "A record" that will be updated by this script.
  - Go to [Your Domains](https://www.hover.com/control_panel/domains) page
  - select a domain then choose the "domain" tab
  - select "Add A Record"
  - leave the type "A" and enter a hostname such as "home"  This will result in a record "home.\<yourdomain>.\<yourtld>"
  - enter any IP address.  This is the what will be updated by this script
  - after it is added, go to `https://www.hover.com/api/domains/<yourdomain>.<yourtld>/dns`
  - look for the "id" field on the same line as the record "name:" you just added. This is what will be used for the HOVERID environment variable below.  It will look like `dns12345678`
  - choose a Poll time in minutes.  This is how often this script will check your current external IP and compare to Hover.  It defaults to 360 min (6 hours) Set this in the POLLTIME environment variable below

## Usage
```
docker run -d \
    --name=hoverddns \
    -e HOVERUSER=<username> \
    -e HOVERPASS=<password> \
    -e HOVERID=<hover id> \
    -e POLLTIME=360 \
    taoofshawn/hoverddns
```
Or use docker-compose.yml:
```yaml
version: '3'
services:
  hoverddns:
    container_name: hoverddns
    image: taoofshawn/hoverddns
    environment:
      - HOVERUSER=<username>
      - HOVERPASS=<password>
      - HOVERID=<hover id>
      - POLLTIME=<minutes between polling hover>
```
use docker logs to monitor:
```
docker logs -f hoverddns
I0410 17:07:40.229029       1 hoverClient.go:50] attempting authentication with Hover
I0410 17:07:41.061048       1 hoverClient.go:72] authenticated successfully with Hover
I0410 17:07:41.061076       1 hoverClient.go:127] checking current IP address setting at Hover
I0410 17:07:41.061087       1 hoverClient.go:93] connecting to Hover
I0410 17:07:41.354872       1 hoverClient.go:110] request to Hover was successful
I0410 17:07:41.355398       1 hoverClient.go:144] checking current external IP address
I0410 17:07:41.489291       1 main.go:34] hover IP needs to be updated. Hover: x.x.x.x, Actual: y.y.y.y
I0410 17:07:41.489376       1 hoverClient.go:158] updating hover IP with y.y.y.y
I0410 17:07:41.489434       1 hoverClient.go:93] connecting to Hover
I0410 17:07:41.546154       1 hoverClient.go:110] request to Hover was successful
I0410 17:07:41.546194       1 main.go:50] sleeping for 360 minutes
```
