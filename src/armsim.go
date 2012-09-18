package main

import (
	"armsim"
	"errors"
	"flag"
	"fmt"
	"log"
)

type Options struct {
	fileName   string
	memorySize uint
	tracing    bool
}

func main() {
	// Setup Logging
	log.SetPrefix("Main: ")

	// Welcome the user
	fmt.Println("ARMSim by Luke Seelenbinder.")

	// Handle command line flags
	options, err := processFlags()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	// Initialize Computer
	c := armsim.NewComputer(uint32(options.memorySize))

	// Disable or Enable tracing
	if !options.tracing {
		c.DisableTracing()
	} else {
		defer c.DisableTracing()
	}

	// Load ELF File
	err = c.LoadELF(options.fileName)
	if err != nil {
		fmt.Println("Unable to load file. Encountered error -", err)
	} else {
		fmt.Println("Loaded valid ELF file - checksum is", c.Checksum())
	}

	// Run the program
	c.Run()
}

func processFlags() (options *Options, err error) {
	// Create Options
	options = new(Options)

	// Define Options
	flag.UintVar(&options.memorySize, "mem", 32768, "RAM size in bytes (1MB max)")
	flag.StringVar(&options.fileName, "load", "", "ELF File Name")
	flag.BoolVar(&options.tracing, "trace", true, "Output trace.log file (default=enabled)")

	// Parse Options
	flag.Parse()

	// Validate Options
	log.Println("RAM Size:", options.memorySize)
	if options.memorySize > 1048576 {
		err = errors.New("RAM size is too large. Must be under 1MB (1048576).")
		return
	}

	log.Println("File name:", options.fileName)
	if options.fileName == "" {
		err = errors.New("Please specify a file name.")
		return
	}

	return
}
