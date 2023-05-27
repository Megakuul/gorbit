package handler

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"github/megakuul/gorbit/logger"
	"io"
	"net"
)

func HandleConnection(srcCon *net.TCPConn, config conf.Config) {
	defer srcCon.Close()

	// Fetching the Endpoint
	index, err := SelectEndpoint(&config.Endpoints)
	if err != nil {
		logger.WriteWarningLogger(err)
		return
	}

	// Append a Session to the Endpoint
	config.Endpoints[index].MutAppendSession()
	defer config.Endpoints[index].MutRemoveSession()
	fmt.Printf("Added a number, current res: %v", config.Endpoints[index].Sessions)

	// Create a TCP Address for the desired Endpoint
	addr, err := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%v", config.Endpoints[index].Hostname, config.Endpoints[index].Port),
	)
	if err != nil {
		logger.WriteInformationLogger("%v", err)
	}

	// Create a TCP Connection to the desired Endpoint
	dstCon, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.WriteInformationLogger(
			"%s failed to reach %s", srcCon.LocalAddr().String(), dstCon.LocalAddr().String(),
		)
		return
	}

	defer dstCon.Close()

	// Redirect dstCon Writer to srcCon Reader
	go func() {
		if err := RedirectBuffered(dstCon, srcCon, config.BufferSizeKB); err != nil {
			logger.WriteInformationLogger("%s", err)
		}
	}()

	// Redirect srcCon Writer to dstCon Reader
	if err := RedirectBuffered(srcCon, dstCon, config.BufferSizeKB); err != nil {
		logger.WriteInformationLogger("%s", err)
	}
}

func RedirectBuffered(dst io.Writer, src io.Reader, size_kb int) error {
	// If the Writer Supports the ReaderFrom interface, it will use this
	// ReadFrom gives it a massiv performance boost
	if reader, ok := dst.(io.ReaderFrom); ok {
		_, err := reader.ReadFrom(src)
		return err
	}

	// If the Reader Supports the WriterTo interface, it will use this
	// WriteTo gives it a massiv performance boost
	if writer, ok := src.(io.WriterTo); ok {
		_, err := writer.WriteTo(dst)
		return err
	}

	// Initialize buffer
	buffer := make([]byte, size_kb*1024)

	// Main Read/Write Loop
	for {
		bytesRead, er := src.Read(buffer)
		if bytesRead > 0 {
			bytesWritten, ew := dst.Write(buffer[0:bytesRead])

			// Catch illegal behaviour
			if bytesWritten < 0 || bytesWritten > bytesRead {
				bytesWritten = 0
				if ew == nil {
					ew = fmt.Errorf("invalid write operation")
				}
			}
			if ew != nil {
				return ew
			}
			if bytesRead != bytesWritten {
				return fmt.Errorf(
					"inconsistent transmission:\n%v bytes read and %v written",
					bytesRead,
					bytesWritten,
				)
			}
		}
		if er != nil {
			if er != io.EOF {
				return er
			}
			return nil
		}
	}
}
