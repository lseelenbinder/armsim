package main

import (
  "armsim/ram"
  "armsim/loader"
  "flag"
  "fmt"
  "log"
  "os"
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
  memory := ram.Init(uint32(options.memorySize))

  // Load ELF File
  success := loader.LoadELF(options.fileName, &memory)
  if success {
    fmt.Println("Loaded valid ELF file - checksum is", memory.Checksum())
  } else {
    fmt.Println("Unable to load file. It may not be an ELF file or it may not exist.")
  }
}

func processFlags() *Options {
  defer func() {
    err := recover()
    if err != nil {
      Usage()
    }
  }()

  // Create Options
  options := new(Options)

  // Define Options
  flag.UintVar(&options.memorySize, "mem", 32768, "RAM size in bytes")
  flag.StringVar(&options.fileName, "load", "", "ELF File Name")

  // Parse Options
  flag.Parse()

  // Validate Options
  log.Println("RAM Size:", options.memorySize)
  if options.memorySize > 1048576 {
    fmt.Println("RAM size is too large. Must be under 1MB (1048576).")
    log.Panicln("RAM size is too large. Quiting...")
  }

  log.Println("File name:", options.fileName)
  if options.fileName == "" {
    fmt.Println("Please specify a file name.")
    log.Panicln("Must specify a file name. Quiting...")
  }

  return options
}

func Usage() {
  fmt.Println("usage: armsim [ --load elf-file ] [ --mem memory-size ]")
  os.Exit(1)
}
