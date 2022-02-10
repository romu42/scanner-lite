package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

func main() {
	// zones is the list of zones the scanner will be monitoring
	zones := ReadConf()
	var child_nses []string

	for zone, parent := range zones {
		parent.child_ns = make(map[string]*Child)
		log.Printf("Working with zone: %s ", zone)

		// Get DS records for zone from Parent
		parent.ds = GetDS(zone, parent.hostname, parent.ip, parent.port)
		for _, ds := range parent.ds {
			log.Printf("%s", ds)
		}

		// get child NSes from Parent and create Child struct
		child_nses = GetNS(zone, parent.hostname, parent.ip, parent.port)
		for _, ns := range child_nses {
			log.Printf("Got NS: %s", ns)
			ip := GetIP(ns, parent.ip, parent.port)
			child := &Child{
				hostname: ns,
				ip:       ip,
				port:     "53",
			}
			parent.child_ns[ns] = child
			log.Printf("%s has ip %s", child.hostname, child.ip)

		}
	}

	// Get Child information
	for zone, parent := range zones {
		log.Printf("Working with Zone: %s", zone)
		for _, child := range parent.child_ns {
			log.Printf("Working with Child NS: %s", child.hostname)
			child.nses = make(map[string]string)

			// Get CDS From Child
			child.cds = GetCDS(zone, child.hostname, child.ip, child.port)
			for _, cds := range child.cds {
				log.Printf("%s", cds)
			}
			// Get CSYNC From Child
			child.csync = GetCsync(zone, child.hostname, child.ip, child.port)
			log.Printf("CSYNC from child: %s", child.csync)

			// Get NSes from Child
			nses := GetNS(zone, child.hostname, child.ip, child.port)
			for _, ns := range nses {
				ip := GetIP(ns, child.ip, child.port)
				log.Printf("IP from child: %s", ip)
				child.nses[ns] = ""
				log.Printf("NS from child: %s", ns)
			}
		}
	}

	// Get DS update information
	for zone, parent := range zones {
		log.Printf("Zone %s\n", zone)
		dsadd, dsremove := CreateUpdate(zone, parent)
		for _, value := range dsadd {
			value.Hdr.Rrtype = 43
		}

		adds := []dns.RR{}
		for _, value := range dsadd {
			adds = append(adds, &value.DS)
		}

		removes := []dns.RR{}
		for _, value := range dsremove {
			removes = append(removes, value)
		}
		log.Printf("value is a %T with value of %v", removes, removes)

		// trying to get ddns to work
		output := []string{}
		args := []string{parent.ip + ":" + parent.port, "catch22.se.", parent.keyname}
		log.Printf("%v", args)

		// nsupdater_updater.go
		updater := GetUpdater("nsupdate")
		err := updater.Update(zone, parent.ip+":"+parent.port, &[][]dns.RR{adds}, &[][]dns.RR{removes}, &output)
		if err != nil {
			fmt.Printf("bob Got an err %v\n", err)
		}
		fmt.Println(output)
	}

	// if CDS's from children match
	// - Update Parent if neccessary ( need to figure out the ttl bit )
	// else
	// - log error

	// Get Csync from children
	// if Csync from children match
	//  - check intent update parent as necessary
	// else
	//  - log error
}

// Ignore CDNSKEY for now
