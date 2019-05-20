package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"net"
	"os"
)

const ttl = int64(60)
const weight = int64(1)

var hostname string
var zoneId string
var domain string

func init() {
	defaultHostname, _ := os.Hostname()
	flag.StringVar(&domain, "domain", "", "domain name")
	flag.StringVar(&hostname, "host", defaultHostname, "hostname for A record (defaults to hostname)")
	flag.StringVar(&zoneId, "zoneid", "", "AWS Zone Id for domain")
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// Get preferred outbound ip of this machine
// Thanks go to Marcel Molina via
// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func createA(svc *route53.Route53, zoneId string, hostname string, ip string) error {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(hostname),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL:           aws.Int64(ttl),
						Weight:        aws.Int64(weight),
						SetIdentifier: aws.String(fmt.Sprintf("A record UPSERT for %s -> %s", hostname, ip)),
					},
				},
			},
			Comment: aws.String(fmt.Sprintf("A record UPSERT for %s -> %s", hostname, ip)),
		},
		HostedZoneId: aws.String(zoneId),
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	// Pretty-print the response data.
	fmt.Println("Change Response:")
	fmt.Println(resp)

	if err != nil {
		return err
	}

	return nil

}

func main() {

	flag.Parse()

	sess, err := session.NewSession()
	if err != nil {
		fatal(err)
	}

	ip := GetOutboundIP()
	fmt.Printf("ip: %#v\n", ip.String())

	svc := route53.New(sess)
	fqdn := fmt.Sprintf("%s.%s", hostname, domain)
	fmt.Printf("Updating for %s -> %s", fqdn, ip)
	err = createA(svc, zoneId, fqdn, ip.String())
	if err != nil {
		fatal(err)
	}

}
