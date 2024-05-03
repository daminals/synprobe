package main

import (
	"testing"
)


// tests
func evalResponse(t *testing.T, response PortResponse, expectedResponse ExpectedResponse) {
	if response.PortType != expectedResponse.PortType {
		t.Errorf("PortType: expected %s, got %s", expectedResponse.PortType, response.PortType)
	}
	if response.ServerInitiated != expectedResponse.ServerInitiated {
		t.Errorf("ServerInitiated: expected %t, got %t", expectedResponse.ServerInitiated, response.ServerInitiated)
	}
	if response.TCP != expectedResponse.TCP {
		t.Errorf("TCP: expected %t, got %t", expectedResponse.TCP, response.TCP)
	}
	if response.TLS != expectedResponse.TLS {
		t.Errorf("TLS: expected %t, got %t", expectedResponse.TLS, response.TLS)
	}
	if response.HTTP != expectedResponse.HTTP {
		t.Errorf("HTTP: expected %t, got %t", expectedResponse.HTTP, response.HTTP)
	}
	if response.HTTPS != expectedResponse.HTTPS {
		t.Errorf("HTTPS: expected %t, got %t", expectedResponse.HTTPS, response.HTTPS)
	}
}

func runResponse(service, port string) PortResponse {
	response := PortResponse{Port: port, Target: service, Response: nil, ServerInitiated: false, Open: false, TCP: false, TLS: false, HTTP: false, HTTPS: false, PortType: "0"}
	response = queryPortTLS(response)
	if !response.TLS {
		response = queryPort(response)
	}
	response.setPortType()
	return response
}

type ExpectedResponse struct {
	PortType string
	ServerInitiated bool
	TCP bool
	TLS bool
	HTTP bool
	HTTPS bool
}

func createTCPServerInitiatedResponse() ExpectedResponse {
	return ExpectedResponse{PortType: "1", ServerInitiated: true, TCP: true, TLS: false, HTTP: false, HTTPS: false}
}

func createTLSServerInitiatedResponse() ExpectedResponse {
	return ExpectedResponse{PortType: "2", ServerInitiated: true, TCP: false, TLS: true, HTTP: false, HTTPS: false}
}

func createTCPClientInitiatedResponse() ExpectedResponse {
	return ExpectedResponse{PortType: "5", ServerInitiated: false, TCP: true, TLS: false, HTTP: false, HTTPS: false}
}

func createTLSClientInitiatedResponse() ExpectedResponse {
	return ExpectedResponse{PortType: "6", ServerInitiated: false, TCP: false, TLS: true, HTTP: false, HTTPS: false}
}

func createHTTPClientInitiatedResponse() ExpectedResponse {
	return ExpectedResponse{PortType: "3", ServerInitiated: false, TCP: true, TLS: false, HTTP: true, HTTPS: false}
}

func createHTTPSClientInitiatedResponse() ExpectedResponse {
	return ExpectedResponse{PortType: "4", ServerInitiated: false, TCP: false, TLS: true, HTTP: false, HTTPS: true}
}

// ftp.dlptest.com 21
func TestDLPTest(t *testing.T) {
	// TCP SERVER INITIATED
	service := "ftp.dlptest.com"
	port := "21"

	response := runResponse(service, port)
	expectedResponse := createTCPServerInitiatedResponse()

	evalResponse(t, response, expectedResponse)
}

// TCP CLIENT INITIATED [HTTP]
// www.cs.stonybrook.edu 80
func TestStonyBrook(t *testing.T) {
	service := "www.cs.stonybrook.edu"
	port := "80"

	response := runResponse(service, port)
	expectedResponse := createHTTPClientInitiatedResponse()
	evalResponse(t, response, expectedResponse)
}

// TLS CLIENT INITIATED [HTTP]
// www.cs.stonybrook.edu 80
func TestStonyBrookTLS(t *testing.T) {
	service := "www.cs.stonybrook.edu"
	port := "443"

	response := runResponse(service, port)
	expectedResponse := createHTTPSClientInitiatedResponse()
	evalResponse(t, response, expectedResponse)
}

// TLS CLIENT INITIATED [HTTPS]
// www.cs.stonybrook.edu 443
func TestGoogle(t *testing.T) {
	service := "www.google.com"
	port := "443"

	response := runResponse(service, port)
	expectedResponse := createHTTPSClientInitiatedResponse()
	evalResponse(t, response, expectedResponse)
}

// TLS GENERIC CLIENT INITIATED
// 8.8.8.8:853
func TestGoogleDNS_DOH(t *testing.T) {
	service := "dns.google.com"
	port := "853"

	response := runResponse(service, port)
	expectedResponse := createTLSClientInitiatedResponse()
	evalResponse(t, response, expectedResponse)
}
// TCP Client Initiated
// cloudflare dns one.one.one.one:53
func TestGoogleDNS(t *testing.T) {
	service := "dns.google.com"
	port := "53"

	response := runResponse(service, port)
	expectedResponse := createTCPClientInitiatedResponse()
	evalResponse(t, response, expectedResponse)
}


// TLS server-initiated
// imap.gmail.com:993
func TestGmailTLS(t *testing.T) {
	service := "imap.gmail.com"
	port := "993"

	response := runResponse(service, port)
	expectedResponse := createTLSServerInitiatedResponse()
	evalResponse(t, response, expectedResponse)
}