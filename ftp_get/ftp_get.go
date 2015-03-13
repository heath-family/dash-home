package ftp_get

import (
	"bytes"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"

	"code.google.com/p/ftp4go"
)

type fields struct {
	host     string
	port     int
	path     string
	user     string
	password string
}

// Returns an error if the port is not an integer
func splitUrl(u *url.URL) (f fields, e error) {
	// Path
	f.path = strings.TrimLeft(u.Path, "/") // Remove leading slash

	// Host and port
	splitHost, splitPort, err := net.SplitHostPort(u.Host)
	if err != nil {
		// could not splithostport
		f.host = u.Host
		f.port = ftp4go.DefaultFtpPort
	} else {
		var p int64
		p, e = strconv.ParseInt(splitPort, 10, 32)
		f.port = int(p)
		f.host = splitHost
	}

	// User/password
	if usr := u.User; usr != nil {
		f.user = usr.Username()
		f.password, _ = usr.Password()
	} else {
		f.user = "anonymous"
	}

	return
}

func Get(u *url.URL) (io.Reader, error) {
	b := bytes.Buffer{}
	ftpClient := ftp4go.NewFTP(0)

	f, err := splitUrl(u)
	if err != nil {
		return nil, err
	}

	_, err = ftpClient.Connect(f.host, f.port, "")
	if err != nil {
		return nil, err
	}

	defer ftpClient.Quit()

	_, err = ftpClient.Login(f.user, f.password, "")
	if err != nil {
		return nil, err
	}

	err = ftpClient.GetBytes(
		ftp4go.RETR_FTP_CMD,
		&b,
		ftp4go.BLOCK_SIZE,
		f.path,
	)
	if err != nil {
		return nil, err
	}

	return &b, nil
}
