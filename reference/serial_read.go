package serialread

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

const start = byte(0x68)

// const end = byte(0x20)

type streamer struct {
	serialPort io.ReadWriteCloser
	netPort    net.Conn
	isTCP      bool
}

func main() {
	streamer := streamer{}
	streamer.Init("/dev/cu.usbserial-1420")
	defer streamer.Close()

	b := make([]byte, 8)
	// previousByte := make([]byte, 1)
	data := make([]byte, 0)

	// bitCounter := 0
	for true {
		streamer.Read(b)
		// copy(previousByte, b)
		// fmt.Printf("%#v\n", b)
		// if bitCounter == 12 {
		// 	fmt.Printf("%#v\n", data)
		// 	data = make([]byte, 0)
		// 	bitCounter = 0
		// 	data = append(data, b[0])
		// } else {
		// 	data = append(data, b[0])
		// 	bitCounter++
		// }
		if b[0] == start {
			fmt.Printf("%08b\n", data)
			// fmt.Printf("%#v\n", data)
			total := 0
			for _, datum := range data {
				total += datum
			}
			data = make([]byte, 0)
			data = append(data, b[0])
		} else {
			data = append(data, b[0])
		}

	}
}

func (s *streamer) Init(addr string) error {

	//Configure the serial port
	options := serial.OpenOptions{
		PortName:              addr,
		BaudRate:              19200,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
		ParityMode:            serial.PARITY_NONE,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		return err
	}

	s.serialPort = port
	s.isTCP = false

	return nil

}

func (s *streamer) Write(p []byte) error {

	if s.isTCP {
		s.netPort.SetReadDeadline(time.Now().Add(1 * time.Second))
		_, err := s.netPort.Write(p)
		if err != nil {
			return err
		}
	} else {
		_, err := s.serialPort.Write(p)
		if err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (s *streamer) Read(p []byte) error {

	_, err := s.serialPort.Read(p)
	if err != nil {
		return err
	}
	return nil

}

func (s *streamer) Close() {
	if s.isTCP {
		s.netPort.Close()
	} else {
		s.serialPort.Close()
	}
}
