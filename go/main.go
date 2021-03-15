package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/spnego"
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
	err := cl.Login()
	errCheck(err)

	r, err := http.NewRequest("GET", url, nil)
	errCheck(err)

	spnegoCl := spnego.NewClient(cl, nil, spn)
	resp, err := spnegoCl.Do(r)
	errCheck(err)

	b, err := ioutil.ReadAll(resp.Body)
	errCheck(err)

	fmt.Println(string(b))
}

// to make ccache file on Mac,
// kinit -c hoge.ccache scott
type Options struct {
	UseKeytab  bool
	KeytabPath string
	UseCcache  bool
	CcachePath string
}

var options = Options{
	UseKeytab:  false,
	KeytabPath: "",
	UseCcache:  true,
	CcachePath: "hoge.ccache",
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func initFlags() {
	flag.BoolVar(&options.UseKeytab, "kt", options.UseKeytab, "Use Keytab")
	flag.BoolVar(&options.UseCcache, "cc", options.UseCcache, "Use Credential cache. Generate it by kinit -c hoge.ccache scott")
	flag.StringVar(&options.KeytabPath, "ktpath", options.KeytabPath, "Keytab path")
	flag.StringVar(&options.CcachePath, "ccpath", options.CcachePath, "Credential cache file path")
	flag.Parse()
}

func main() {
	var err error
	var kt *keytab.Keytab
	var ccache *credentials.CCache
	var c *config.Config
	var cl *client.Client

	initFlags()

	// Enable below if you want to read /etc/krb5.conf
	// Note:
	// MacOS uses Heimdal and krb5.conf format is different (e.g. tcp/$hostname)
	// To use gokrb, it must be MIT format (only $hostname)
	// that is the reason why I keep the conf in KRB5CONF string const above
	//
	// krb5ConfReader, err := os.Open("/etc/krb5.conf")
	// errCheck(err)
	// defer krb5ConfReader.Close()
	// c, err = config.NewFromReader(krb5ConfReader)

	c, err = config.NewFromString(KRB5CONF)
	errCheck(err)
	c.LibDefaults.NoAddresses = true

	if options.UseKeytab {
		kt, err = keytab.Load(options.KeytabPath)
		errCheck(err)
		user := os.Getenv("USER")
		realm := "REALM.SOKOIDE.COM"
		cl = client.NewWithKeytab(user, realm, kt, c)
	} else if options.UseCcache {
		ccache, err = credentials.LoadCCache(options.CcachePath)
		errCheck(err)
		cl, err = client.NewFromCCache(ccache, c)
	} else {
		panic("You must use either Keytab or Ccache")
	}

	httpRequest("http://nginx-spnego:20080", "HTTP/nginx-spnego", cl)
}
