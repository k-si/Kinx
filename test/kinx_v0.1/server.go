package main

import "kinx/knet"

func main() {
	s := knet.NewServer("[kinx v0.1]")
	s.Serve()
}
