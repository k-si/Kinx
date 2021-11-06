package main

import "kinx/knet"

func main() {
	s := knet.NewServer("Kinx v0.2")
	s.Serve()
}
