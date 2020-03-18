# hoverDdns
Keep a hover DNS record up to date with your external IP. 
A use-case would be a home LAN with the following features:
- NAT router - with a public WAN ip that may change periodically
- Some NAS device or other "always on" computer capable of running Docker
- Some desire to always know the external IP (i.e. VPN)
- you already have a domain on Hover for some other purpose

I have been using one of the free dynamic dns providers with a monthly nag email and didnt hear about duckdns until i got too far and interested in finishing this one :/

Info | Based on [hover-dns-updater.](https://github.com/texasaggie97/hover-dns-updater) Simplified and targetted for a simple docker image


