armsim
======

Overview
--------

Simulation of the ARM architecture written in Go. At the present,
only the memory and loader functionality are completed.

Requirements
------------

- [Go](http://www.golang.org/)

Build and Test
-------------

Go can be complicated to build the first time. The general methodology is to
obtain the version of Go recommended for your platform
[http://golang.org/doc/install](http://golang.org/doc/install). Don't forget to
set a proper `$GOPATH`.

After obtaining Go, build the project on Linux by:
- `cd` to the armsim/ folder in the project.
- Add the armsim/ folder to your GOPATH: `export GOPATH=$GOPATH:`pwd``
- Build the armsim executable with: `go build -o install/armsim src/armsim.go`
- Change to the install directory and run `./armsim` (with `2> /dev/null` to avoid
  log entries)

To test, run: `go test armsim/loader armsim/ram`

Configuration
-------------

Logging at this point does not have a flag to turn it off; however, all the
output is to STDERR, so you can use redirection to "turn it off"; see build
section.

User Guide
---------

To run the project, simply `cd` to the install/ directory of the project and
run the armsim executable.
