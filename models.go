package main

import (
	"log"
	"os"
)

var lg *log.Logger = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lmicroseconds)

type Profile struct{
  Name string
  Age uint8
  Gender string
  Interest string
  Description string
  Photo string
}
