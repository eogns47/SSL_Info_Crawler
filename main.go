package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func extractOrganization(issuer string) string {
	parts := strings.Split(issuer, ",")
	for _, part := range parts {
		if strings.HasPrefix(part, "O=") {
			return strings.TrimPrefix(part, "O=")
		}
	}
	return ""
}

func getTLSProtocol(conn *tls.Conn) string {
	version := conn.ConnectionState().Version
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

func connectServer(target string, port string) ([]*x509.Certificate, string, error) {
	urlWithPort := target + ":" + port
	conn, err := tls.Dial("tcp", urlWithPort, &tls.Config{})
	if err != nil {
		errors.Wrap(err, "Error connecting to server: "+target)
		return nil, "", err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	protocol := getTLSProtocol(conn)
	return certs, protocol, nil
}

func getSubject(url string) (string, error) {
	parts := strings.Split(url, ".")
	if len(parts) < 2 {
		return "", errors.New("Invalid url")
	}
	subject := parts[len(parts)-2] + "." + parts[len(parts)-1]

	return subject, nil
}

func main() {
	//readCsv()
	//길이가 2보다 작으면 사용법 출력

	if len(os.Args) < 3 {
		fmt.Println("😅Usage:     go run main.go <url> <port> \n🤔Example:   go run main.go www.google.com 443")
		return
	}

	target := os.Args[1]
	port := os.Args[2]
	certs, protocol, err := connectServer(target, port)

	if err != nil {
		fmt.Println("⚠️Error with url: ", err)
		return
	}
	// 서버의 인증서 가져오기
	if len(certs) == 0 {
		fmt.Println("No certificates found.")
		return
	}

	// 첫 번째 인증서 선택
	cert := certs[0]

	// 발급자 정보 출력
	issuer := cert.Issuer.String()
	organization := extractOrganization(issuer)
	preDate := cert.NotBefore.Format("2006-01-02")
	postDate := cert.NotAfter.Format("2006-01-02")

	certificateLife_Year := int(cert.NotAfter.Sub(cert.NotBefore).Hours() / 24 / 365)
	certificateLife_Month := int(cert.NotAfter.Sub(cert.NotBefore).Hours()/24) % 365 / 30
	certificateLife_Day := int(cert.NotAfter.Sub(cert.NotBefore).Hours()/24) % 365 % 30

	certificateLife := fmt.Sprintf("%d years %d months %d days", certificateLife_Year, certificateLife_Month, certificateLife_Day)

	subject, nil := getSubject(target)
	if err != nil {
		fmt.Println("⚠️Error with url: ", err)
		return
	}

	fmt.Println("URL                   : ", target)
	fmt.Println("Subject               : ", subject)
	fmt.Println("Issuer Organization   : ", organization)
	fmt.Println("valid_since           : ", preDate)
	fmt.Println("valid_until           : ", postDate)
	fmt.Println("Certificate Life      : ", certificateLife)
	fmt.Println("TLS Protocol          : ", protocol)
}
