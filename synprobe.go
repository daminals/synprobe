package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)


var PROBES = map[string]string{"http": "GET / HTTP/1.0\r\n\r\n", "generic": "\r\n\r\n\r\n\r\n"}

// read CA certs
func readCerts() *tls.Config {
	return &tls.Config{
		// RootCAs: roots,
	}
}

type PortResponse struct {
	Target string
	Port string
	Response []byte
	ServerInitiated bool
	Open bool
	TCP bool
	TLS bool
	HTTP bool
	HTTPS bool
	PortType string
}

func (p *PortResponse) setPortType() {
	if p.HTTPS {
		// HTTPS
		p.PortType = "4"
	} else if p.HTTP {
		// HTTP
		p.PortType = "3"
	} else if p.TLS {
		if p.ServerInitiated {
			// TLS SERVER INITIATED
			p.PortType = "2"
		} else {
			// TLS CLIENT INITIATED
			p.PortType = "6"
		}
	} else if p.TCP {
		if p.ServerInitiated {
			// TCP SERVER INITIATED
			p.PortType = "1"
		} else {
			// TCP CLIENT INITIATED
			p.PortType = "5"
		}
	} else {
		// idk
		p.PortType = "0"
	}
}


// portResponse toString
func (p PortResponse) String() string {
	if !p.Open {
		return "Port: " + p.Port + "\nClosed" 
	}

	// now need to determine what type of connection it is
	// start from most specific to least specific
	p.setPortType()
	str := "Port: " + p.Port + "\n" + "Port Type: " + p.PortType + "\nServerInitiated: " + strconv.FormatBool(p.ServerInitiated) + "\n" + "TCP: " + strconv.FormatBool(p.TCP) + "\n" + "TLS: " + strconv.FormatBool(p.TLS) + "\n" + "HTTP: " + strconv.FormatBool(p.HTTP) + "\n" + "HTTPS: " + strconv.FormatBool(p.HTTPS)
	str += "\n----------------------------------\n" + formatBytes(p.Response) + "\n----------------------------------\n"
	return str
}

func formatBytes(buf []byte) string {
	response := ""
	// check if the byte is a printable character
	for _, b := range buf {
		if b <= 126 {
			response += string(b)
		} else {
			response += "."
		}
	}
	return response
}

// The range of ports to be scanned (just a single number for one port,
//	or a port range in the form X-Y for multiple ports).
func decodePortRange(portRange string) (ports []string) {
	if portRange == "" {
		// By default, if '-p' is not provided, the tool should scan only for the
		// following commonly used TCP ports: 21, 22, 23, 25, 80, 110, 143, 443, 587,
		// 853, 993, 3389, 8080.
		return []string{"21", "22", "23", "25", "80", "110", "143", "443", "587", "853", "993", "3389", "8080"}
	}

	if strings.Contains(portRange, "-") {
		// port range
		// split port range
		portStrings := strings.Split(portRange, "-")
		if len(portStrings) != 2 {
			log.Fatal("Invalid port range")
		}
		// convert to int
		start, err := strconv.Atoi(portStrings[0])
		if err != nil {
			log.Fatal("Invalid port range")
		}
		end, err := strconv.Atoi(portStrings[1])
		if err != nil {
			log.Fatal("Invalid port range")
		}
		// check range
		if start > end {
			log.Fatal("Invalid port range")
		}
		ports = make([]string, 0)
		// create port range
		for i := start; i <= end; i++ {
			ports = append(ports, strconv.Itoa(i))
		}
		return ports
	}

	// single port is the only option left
	if _, err := strconv.Atoi(portRange); err == nil {
		// single port
		return []string{portRange}
	} 

	log.Fatal("Invalid port range")
	return nil
}

func readResponse(conn net.Conn) (int, []byte) {
	buf := make([]byte, 1024)
	// set a 3 second timeout for the connection
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	// read until 1024 or EOF
	total := 0
	for total < 1024 {
		// read from the connection
		localBuf := make([]byte, 1024)
		n, err := conn.Read(localBuf)
		if err != nil {
			break
		}
		total += n
		// append to the buffer, but only up to 1024 bytes
		if total > 1024 {
			excess := total - 1024
			localBuf = localBuf[:n-excess]
		}
		buf = append(buf, localBuf...)
	}
	return total, buf
}


func queryPort(response PortResponse) PortResponse {
		if response.TLS {
			return response
		}


		// make a tcp connection to the target
		conn, err := net.Dial("tcp", response.Target + ":" + response.Port)
		if err != nil {
			// log.Println("Error connecting to", target + ":" + port)
			response.Open = false
			return response
		}
		response.Open = true
		response.TCP = true

		// EXPECT SERVER RESPONSE SECTION
		// print the first 1024 bytes returned by the server (server-initiated dialog)
		total, buf := readResponse(conn)
		// check if the response is empty
		if total != 0 {
			response.ServerInitiated = true
			// set the response
			response.Response = buf
			return response
		}
		conn.Close()


	// in case the
	// server doesn't send any data after 3 seconds, try to elicit a response by
	// sending a series of probe requests (and if a probe request succeeds, again
	// print the first 1024 bytes returned)
	 
	for name, probe := range PROBES {
		conn, err := net.Dial("tcp", response.Target + ":" + response.Port)
		if err != nil {
			// log.Println("Error connecting to", response.Target + ":" + response.Port)
			return response
		}
		defer conn.Close()
		// send the probe
		_, err = conn.Write([]byte(probe))
		if err != nil {
			// log.Println("Error sending probe to", response.Target + ":" + response.Port)
			return response
		}
		// print the first 1024 bytes returned by the server
		n, buf := readResponse(conn)

		// fmt.Println(total, string(buf))

		// set the response
		// check if the response is empty
		if n <= 0 {
			continue
		}

		response.Response = buf
		// check probe type
		switch name {
			case "http":
				response.HTTP = true
			case "generic":
				// pretend to do something
			default:
				log.Fatal("Invalid probe type")
		}
	}


	return response
}

func queryPortTLS(response PortResponse) PortResponse {
		conf := readCerts()
		// make a tls connection to the target
		conn, err := tls.Dial("tcp", response.Target + ":" + response.Port, conf)

		if err != nil {
			// log.Println("Error connecting to", response.Target + ":" + response.Port)
			response.Open = false
			return response
		}
		defer conn.Close()
		response.Open = true
		response.TLS = true

		// EXPECT SERVER RESPONSE SECTION
		// print the first 1024 bytes returned by the server (server-initiated dialog)
		total, buf := readResponse(conn)
		// check if the response is empty
		if total != 0 {
			response.ServerInitiated = true
			// set the response
			response.Response = buf
			return response
		}

	// in case the server doesn't send any data after 3 seconds, try to elicit a response by
	// sending a series of probe requests (and if a probe request succeeds, again print the first 1024 bytes returned)
	for name, probe := range PROBES {
		// send the probe
		_, err := conn.Write([]byte(probe))
		if err != nil {
			// log.Println("Error sending probe to", response.Target + ":" + response.Port)
			return response
		}
		// print the first 1024 bytes returned by the server
		total, buf := readResponse(conn)

		// set the response
		// check if the response is empty
		if total == 0 {
			continue
		}

		response.Response = buf
		// check probe type
		switch name {
			case "http":
				if strings.Contains(string(buf), "HTTP") {
					response.HTTPS = true
				}
			case "generic":
				// pretend to do something
			default:
				log.Fatal("Invalid probe type")
		}
	}

	return response
}



func main() {
		// go run synprobe.go [ -p portRange ] target
		portRange := flag.String("p", "", "port_range")
		flag.Parse()
	
		// check arg length
		if len(flag.Args()) != 1 {
			log.Fatal("Usage: synprobe [ -p port_range ] target")
		}

		ports := decodePortRange(*portRange)
		target := flag.Args()[0]

		
		// iterate over the ports
		for _, port := range ports {
			response := PortResponse{Port: port, Target: target, Response: nil, ServerInitiated: false, Open: false, TCP: false, TLS: false, HTTP: false, HTTPS: false, PortType: "0"}
			// check if the port is open
			response = queryPortTLS(response)
			if !response.TLS {
				response = queryPort(response)
			}
			fmt.Println(response.String() + "\n")
		}
}
