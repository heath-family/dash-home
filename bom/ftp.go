package bom

import (
	"code.google.com/p/ftp4go"
	"io"
)

const (
	ftpHost        = "ftp2.bom.gov.au"
	remoteFilename = "anon/gen/fwo/IDV10450.xml"
)

func fetch() (io.Reader, error) {
	reader, writer := io.Pipe()
	ftpClient := ftp4go.NewFTP(0)
	_, err := ftpClient.Connect(ftpHost, ftp4go.DefaultFtpPort, "")
	if err != nil {
		return nil, err
	}

	defer ftpClient.Quit()

	_, err = ftpClient.Login("anonymous", "", "")
	if err != nil {
		return nil, err
	}
	err = ftpClient.GetBytes(
		ftp4go.RETR_FTP_CMD,
		writer,
		ftp4go.BLOCK_SIZE,
		remoteFilename,
	)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
