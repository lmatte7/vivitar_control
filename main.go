package main

import (
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

const start = byte(0x68)

type streamer struct {
	serialPort  io.ReadWriteCloser
	BaseCommand []byte
}

func main() {
	streamer := streamer{}
	// streamer.Init("/dev/cu.usbserial-1420")
	streamer.Init("/dev/serial0")
	defer streamer.Close()

	for i := 0; i < 20; i++ {
		err := streamer.SendUpCommand()
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 50; i++ {
		err := streamer.SendDownCommand()
		if err != nil {
			panic(err)
		}
	}
}

func (s *streamer) SendUpCommand() error {

	upCommand := make([]byte, 12)
	copy(upCommand, s.BaseCommand)

	upCommand[5] = 0xff

	err := s.Write(upCommand)

	if err != nil {
		return err
	}

	return nil
}

func (s *streamer) SendDownCommand() error {

	upCommand := make([]byte, 12)
	copy(upCommand, s.BaseCommand)

	upCommand[5] = 0x00

	err := s.Write(upCommand)

	if err != nil {
		return err
	}

	return nil
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

	s.BaseCommand = []byte{0x68, 0x01, 0x09, 0x80, 0x80, 0x80, 0x80, 0x20, 0x08, 0x00, 0x00, 0x01}
	s.serialPort = port

	return nil

}

func (s *streamer) Write(p []byte) error {

	_, err := s.serialPort.Write(p)
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

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
	s.serialPort.Close()
}
