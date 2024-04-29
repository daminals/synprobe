# Synprobe

## Description
Synprobe is a golang tool that can be used to probe a network for open ports. It is a simple tool that identifies open ports by the service they offer using 6 categories

1) TCP server-initiated (server banner was immediately returned over TCP)
2) TLS server-initiated (server banner was immediately returned over TLS)
3) HTTP server (GET request over TCP successfully elicited a response)
4) HTTPS server (GET request over TLS successfully elicited a response)
5) Generic TCP server (Generic lines over TCP may or may not elicit a response)
6) Generic TLS server (Generic lines over TLS may or may not elicit a response)

