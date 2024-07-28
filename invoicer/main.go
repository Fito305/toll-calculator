package main

import (
	"fmt"

)

func main() {
	fmt.Println("this is working fine")
}



// This invoicer is going to have a transport.
// How are we going to reach this invoicer?
// Are we going to use JSON or ProtoBuffers?
// We are going to use both! 
// Why? You need to set up your Micro Service `transport independant`.
// We first start off with JSON it's easier to debug. And that is how companies start.
// And once everything is set up we can add another transport that is going to 
// be the Protobuffers. Knowing how to implement Protobuffer is very useful.
