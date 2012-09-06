package main

import (
  "armsim/ram"
  "flag"
  "fmt"
  "log"
)

type Options struct {
  fileName   string
  memorySize uint
}

func main() {
  // Setup Logging
  log.SetPrefix("Main: ")

  // Welcome the user
  fmt.Println("ARMSim by Luke Seelenbinder.")

  // Handle command line flags
  options := processFlags()

  // Initialize RAM
  ram.Init(uint32(options.memorySize))
}

func processFlags() *Options {
  // Create Options
  options := new(Options)

  // Define Options
  flag.UintVar(&options.memorySize, "mem-size", 32768, "RAM size in bytes")
  flag.StringVar(&options.fileName, "load", "", "ELF File Name")

  // Parse Options
  flag.Parse()

  // Validate Options
  log.Println("RAM Size:", options.memorySize)
  if options.memorySize > 1048576 {
    fmt.Println("RAM size is too large. Must be under 1MB (1048576).")
    log.Fatalln("RAM size is too large. Quiting...")
  }

  log.Println("File name:", options.fileName)
  if options.fileName == "" {
    fmt.Println("Please specify a file name.")
    log.Fatalln("Must specify a file name. Quiting...")
  }

  return options
}
