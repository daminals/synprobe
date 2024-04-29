# Synprobe
Synprobe is a golang tool that can be used to probe a network for open ports.

## Description
 Synprobe identifies open ports by the service they offer using 6 categories

1) TCP server-initiated (server banner was immediately returned over TCP)
2) TLS server-initiated (server banner was immediately returned over TLS)
3) HTTP server (GET request over TCP successfully elicited a response)
4) HTTPS server (GET request over TLS successfully elicited a response)
5) Generic TCP server (Generic lines over TCP may or may not elicit a response)
6) Generic TLS server (Generic lines over TLS may or may not elicit a response)

Usage:
synprobe [-p port_range] target 

  -p  The range of ports to be scanned (just a single number for one port,
      or a port range in the form X-Y for multiple ports). Default is [21, 22, 23, 25, 80, 110, 143, 443, 587, 853, 993, 3389, 8080]

  target  The target to be scanned. This can be an IP address or a domain name.

## How to run
version: go 1.22.2

1. Clone the repository
2. Run `go build synprobe.go`
3. Run `./synprobe -p 443 www.google.com`

## Tests
To run tests, run `go test`
There are tests for each category available in the synprobe_test.go file

Here is some example output:

```
$ ./synprobe -p 443 www.google.com
Port: 443
Port Type: 4
ServerInitiated: false
TCP: false
TLS: true
HTTP: false
HTTPS: true
----------------------------------
24 02:40:30 GMT; path=/; domain=.google.com; HttpOnly
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
Accept-Ranges: none
Vary: Accept-Encoding

<!doctype html><html itemscope="" itemtype="http://schema.org/WebPage" lang="en"><head><meta content="Search the world's information, including webpages, images, videos and more. Google has many special features to help you find exactly what you're looking for." name="description"><meta content="noodp, " name="robots"><meta content="text/html; charset=UTF-8" http-equiv="Content-Type"><meta content="/images/branding/googleg/1x/googleg_standard_color_128dp.png" itemprop="image"><title>Google</title><script nonce="Ja11vmn3UK17REWLN4OmLg">(function(){var _g={kEI:'nggvZtO9LbXZ5NoPmamEOA',kEXPI:'0,202714,1167764,1133897,1195935,1000,478218,4998,58391,2891,3926,7828,76620,781,1244,1,16916,65253,24058,6636,49751,2,16395,342,23024,6699,41949,24670,33064,2,2,1,24626,2006,8155,8860,14490,22436,5600,4179,12414,26264,3781,20199,36746,3801,2412,30219,3030,15816,1804,7
----------------------------------
```