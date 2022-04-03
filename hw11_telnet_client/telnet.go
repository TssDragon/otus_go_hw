package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrHasNoOpenConnection = errors.New("there is no open connection")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}
	return c.connection.Close()
}

func (c *Client) Send() error {
	if c.connection == nil {
		return ErrHasNoOpenConnection
	}

	_, err := io.Copy(c.connection, c.in)
	return err
}

func (c *Client) Receive() error {
	if c.connection == nil {
		return ErrHasNoOpenConnection
	}

	_, err := io.Copy(c.out, c.connection)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
