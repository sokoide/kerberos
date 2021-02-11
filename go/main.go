package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

const (
	port     = ":9080"
	KRB5CONF = `[libdefaults]
	default_realm = REALM.SOKOIDE.COM
	kdc_timesync = 1
	ccache_type = 4
	forwardable = true
	proxiable = true

	allow_weak_crypto = true
	dns_lookup_realm = false
	dns_lookup_kdc = false
	ticket_lifetime = 24h
	rdns = false

[realms]
	REALM.SOKOIDE.COM = {
			kdc = scottmm.local:10088
			admin_server = scottmm.local:10749
			kpasswd_server = scottmm.local:10464
	}
[domain_realm]
	.realm.sokoide.com = REALM.SOKOIDE.COM
	realm.sokoide.com = REALM.SOKOIDE.COM
`
)

func httpRequest(url string, spn string, cl *client.Client) {
	l := log.New(os.Stderr, "GOKRB5 Client: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := cl.Login()
	if err != nil {
		l.Printf("Error on AS_REQ: %v\n", err)
	}
	r, _ := http.NewRequest("GET", url, nil)
	err = spnego.SetSPNEGOHeader(cl, r, spn)
	if err != nil {
		l.Printf("Error setting client SPNEGO header: %v", err)
	}
	httpResp, err := http.DefaultClient.Do(r)
	if err != nil {
		l.Printf("Request error: %v\n", err)
	}
	fmt.Fprintf(os.Stdout, "Response Code: %v\n", httpResp.StatusCode)
	content, _ := ioutil.ReadAll(httpResp.Body)
	fmt.Fprintf(os.Stdout, "Response Body:\n%s\n", content)
}

func httpRequest2(url string, spn string, cl *client.Client) {
	l := log.New(os.Stderr, "GOKRB5 Client: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := cl.Login()
	if err != nil {
		l.Printf("Error on AS_REQ: %v\n", err)
	}

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		l.Fatalf("could create request: %v", err)
	}

	spnegoCl := spnego.NewClient(cl, nil, spn)
	resp, err := spnegoCl.Do(r)
	if err != nil {
		l.Fatalf("error making request: %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Fatalf("error reading response body: %v", err)
	}
	fmt.Println(string(b))
}

func main() {
	// read keytab
	kt, err := keytab.Load("sokoide.keytab")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	// read krb5.conf
	c, err := config.NewConfigFromString(KRB5CONF)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	c.LibDefaults.NoAddresses = true
	cl := client.NewClientWithKeytab("sokoide", "REALM.SOKOIDE.COM", kt, c)
	// 1. manual SPN set
	// httpRequest("http://nginx-spnego:20080", "HTTP/nginx-spnego" cl)

	// 2. automatic by spnego module
	httpRequest2("http://nginx-spnego:20080", "HTTP/nginx-spnego", cl)
}
