package main

import (
	"armsim"
	"armsim/web"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

type Options struct {
	fileName   string
	memorySize uint
	tracing    bool
	gui        bool
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

	halting := make(chan bool, 1)
	finishing := make(chan bool, 1)

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

	if options.gui {
		log.Println("Loading webserver...")
		fmt.Println("Please open your web browser to http://localhost:4567/ to see the gui.")

		// Attempt to open Firefox
		cmd := exec.Command("firefox", "http://localhost:4567/")
		cmd.Start()

		s := web.Server{c, options.fileName, halting, finishing}
		// Launch the webserver
		s.Launch()
	} else {
		// Run the program
		c.Run(halting, finishing)
	}
}

func processFlags() (options *Options, err error) {
	// Create Options
	options = new(Options)

	// Define Options
	flag.UintVar(&options.memorySize, "mem", 32768, "RAM size in bytes (1MB max)")
	flag.StringVar(&options.fileName, "load", "", "ELF File Name")
	flag.BoolVar(&options.tracing, "trace", true, "Output trace.log file (default=enabled)")
	flag.BoolVar(&options.gui, "gui", true, "Use gui instead of command line")

	// Parse Options
	flag.Parse()

	// Validate Options
	log.Println("RAM Size:", options.memorySize)
	if options.memorySize > 1048576 {
		err = errors.New("RAM size is too large. Must be under 1MB (1048576).")
		return
	}

	if !options.gui {
		log.Println("File name:", options.fileName)
		if options.fileName == "" {
			err = errors.New("Please specify a file name.")
			return
		}
	}

	return
}
