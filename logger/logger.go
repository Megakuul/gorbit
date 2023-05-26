package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

var (
	ERROR       bool
	WARNING     bool
	INFORMATION bool
)

func InitLogger(logname string, loglevel string) error {

	file, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	mOutWriter := io.MultiWriter(os.Stdout, file)

	log.SetOutput(mOutWriter)

	logleverUpper := strings.ToUpper(loglevel)

	ERROR = strings.Contains(logleverUpper, "ERROR")
	WARNING = strings.Contains(logleverUpper, "WARNING")
	INFORMATION = strings.Contains(logleverUpper, "INFORMATION")

	return nil
}

func WriteInformationLogger(message string) {
	if INFORMATION {
		log.Printf("\n[ Gorbit Information ]:\n%s\n\n", message)
	}
}

func WriteWarningLogger(err error) {
	if WARNING {
		log.Printf("\n[ Gorbit Warning ]:\n%s\n\n", err)
	}
}

func WriteErrLogger(err error) {
	if ERROR {
		log.Fatalf("\n[ Gorbit Panic ]:\n%s\n\n", err)
	}
}
