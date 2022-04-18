package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Output format is: domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord \n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// checkDomain will lookup the provided domain & try to fetch it's records
func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v \t", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	textRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v \t", err)
	}

	for _, record := range textRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v \n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=dmarc=1") {
			hasDMARC = true
			dmarcRecord = record

			break
		}
	}

	// print out domain records - if domain is valid
	fmt.Printf("%v, %v, %v, %v, %v, %v \n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
