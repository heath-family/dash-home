package ftp_get

import (
	"io/ioutil"
	"net/url"
	"testing"
)

const (
	ftpHost        = "ftp2.bom.gov.au"
	remoteFilename = "anon/gen/fwo/IDV10450.xml"
)

func TestUrlSplit(t *testing.T) {
	testCases := []struct {
		url string
		fields
	}{
		{"ftp://ftp2.bom.gov.au/anon/gen/fwo/IDV10450.xml", fields{
			host:     "ftp2.bom.gov.au",
			port:     21,
			path:     "anon/gen/fwo/IDV10450.xml",
			user:     "anonymous",
			password: "",
		}}, {
			"ftp://foo:bar@ftp2.bom.gov.au:90///bar.baz", fields{
				host:     "ftp2.bom.gov.au",
				port:     90,
				path:     "bar.baz",
				user:     "foo",
				password: "bar",
			},
		},
	}
	for _, c := range testCases {
		u, err := url.Parse(c.url)
		die(err)
		parts, err := splitUrl(u)
		die(err)
		if parts != c.fields {
			t.Errorf("URL '%s' :\n expected: %+v\n      got: %+v", c.url, parts, c.fields)
		}
	}
}

// Actually goes and connects to iinet
func TestFtpGet(t *testing.T) {
	var err error
	u, err := url.Parse("ftp://ftp.iinet.net.au/notfound.html")
	die(err)
	resp, err := Get(u)
	die(err)
	bytes, err := ioutil.ReadAll(resp)
	die(err)
	expected := "Content not found"
	actual := string(bytes[:17])
	if expected != actual {
		t.Errorf("Expected IInet to have a notfound message on their ftp server.\n Expected '%s'\n      got '%s'", expected, actual)
	}
}

func die(err error) {
	if err != nil {
		die(err)
	}
}
