package main

import (
	"log"
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

		// get child NSes from Parent and create ZoneAuth struct
		child_nses = GetNS(zone, parent.hostname, parent.ip, parent.port)
		for _, ns := range child_nses {
			log.Printf("Got NS: %s", ns)
			ip := GetIP(ns, parent.ip, parent.port)
			z := &Child{
				hostname: ns,
				ip:       ip,
				port:     "53",
			}
			parent.child_ns[ns] = z
			log.Printf("%s has ip %s", z.hostname, z.ip)

			// Get CDS From Child
			//z.cds = GetCDS(zone, z.hostname, z.ip, z.port)
			//for _, cds := range z.cds {
			//log.Printf("%s", cds)
			//}

		}
	}

	// Get
	for zone, parent := range zones {
		log.Printf("Working with Zone: %s", zone)
		for _, ns := range parent.child_ns {
			log.Printf("Working with Child NS: %s", ns.hostname)
			// Get CDS From Child
			//z.cds = GetCDS(zone, z.hostname, z.ip, z.port)
			//for _, cds := range z.cds {
			//log.Printf("%s", cds)
			//}
		}
	}

	// Get NSes from Child
	// Get NS A record from Child if they exist
	// if CDS's from children match
	// - Update Parent if neccessary
	// else
	// - log error

	// Ignore CDNSKEY for now
	// Get Csync from children
	// if Csync from children match
	//  - check intent update parent as necessary
	// else
	//  - log error
}
