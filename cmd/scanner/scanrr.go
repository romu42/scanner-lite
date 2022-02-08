package main

import (
	"log"
	"strings"

	//"github.com/miekg/dns"

	"github.com/miekg/dns"
)

func GetIP(hostname string, server string, port string) string {
	log.Printf("Getting %s IP from parent\n", hostname)
	m := new(dns.Msg)
	m.SetQuestion(hostname, dns.TypeA)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, server+":"+port)
	if err != nil {
		log.Printf("%s: Unable to fetch ip from %s: %s", hostname, server, err)
	}

	// Only grabbing the first RR
	answer := r.Answer[0].String()
	ip := answer[strings.LastIndex(answer, "\t")+1:]
	// debug	fmt.Println(ip)
	return ip
}

func GetNS(zone string, hostname, server string, port string) []string {
	log.Printf("Getting %s NSes from %s %s\n", zone, hostname, server)
	m := new(dns.Msg)
	m.SetQuestion(zone, dns.TypeNS)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, server+":"+port)
	if err != nil {
		log.Printf("%s: Unable to fetch NSes from %s: %s", zone, server, err)
	}

	var nses []string
	for _, a := range r.Ns {
		ns, ok := a.(*dns.NS)
		if !ok {
			continue
		}
		nses = append(nses, ns.String()[strings.LastIndex(ns.String(), "\t")+1:])
	}
	log.Printf("nses: %+v\n", nses)
	return nses
}

func GetDS(zone string, hostname string, server string, port string) []string {
	log.Printf("Getting %s DSes from %s %s\n", zone, hostname, server)
	m := new(dns.Msg)
	m.SetQuestion(zone, dns.TypeDS)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, server+":"+port)
	if err != nil {
		log.Printf("%s: Unable to fetch DSes from %s: %s", zone, server, err)
	}

	var dses []string
	for _, a := range r.Answer {
		ds, ok := a.(*dns.DS)
		if !ok {
			continue
		}
		dses = append(dses, ds.String()[strings.LastIndex(ds.String(), "\t")+1:])
	}
	return dses
}

func GetCDS(zone string, hostname string, server string, port string) []string {
	log.Printf("Getting %s CDSes from %s %s\n", zone, hostname, server)
	m := new(dns.Msg)
	m.SetQuestion(zone, dns.TypeCDS)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, server+":"+port)
	if err != nil {
		log.Printf("%s: Unable to fetch CDSes from %s: %s", zone, server, err)
	}

	var cdses []string
	log.Printf("%+v", r.Answer)
	for _, a := range r.Answer {
		cds, ok := a.(*dns.CDS)
		if !ok {
			continue
		}
		cdses = append(cdses, cds.String()[strings.LastIndex(cds.String(), "\t")+1:])
	}
	return cdses
}

/*
func (ns NameServer) GetCds()     {}
func (ns NameServer) GetCdnskey() {}
func (ns Nameserver) GetCsync()   {}
*/
