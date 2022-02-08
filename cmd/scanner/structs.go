package main

// Parent Server info and child info according to Parent
type ZoneP struct {
	hostname string
	ip       string
	port     string
	hmac     string
	keyname  string
	secret   string
	child_ns map[string]*ZoneAuth
	ds       []string
}

// Authoritave Nameserver
type ZoneAuth struct {
	hostname string
	ip       string
	port     string
	nses     map[string]string
	cds      []string
	cdnskey  []string
	csync    string
}
