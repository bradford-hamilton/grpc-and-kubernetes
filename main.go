package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// Args defines what can be passed to our rpc handler methods
type Args struct {
	A, B int
}

// Quotient defines the shape of what is returned from a Divide rpc method call. It
// includes the Quotient and the remainder
type Quotient struct {
	Quotient, Remainder int
}

// Arithmatic is a type alias for an integer which we use to write methods off of
type Arithmatic int

// Multiply takes args A & B, multiplies them, then assigns result to pointer to the reply
func (t *Arithmatic) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

// Divide takes args A & B, checks that B is not zero, then finds the quotient & remainder
// and assigns them to the *Quotient
func (t *Arithmatic) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("cannot divide by zero")
	}

	quo.Quotient = args.A / args.B
	quo.Remainder = args.A % args.B

	return nil
}

func main() {
	// Create and register our Arithmatic rpc service
	arith := new(Arithmatic)
	rpc.Register(arith)

	// Returns us an address of a TCP endpoint
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1337")
	handleErrors(err)

	// Listen to the tcp endpoint returned from "ResolveTCPAddr"
	listener, err := net.ListenTCP("tcp", tcpAddr)
	handleErrors(err)

	// Loop forever listening for connections to our rpc service
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}
}

// handleErrors simply takes an err and checks whether or not it's nil. If we have an error
// it prints the error and exits the program with a status code of 1
func handleErrors(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
