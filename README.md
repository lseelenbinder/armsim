armsim
======

Overview
--------

Simulation of the ARM architecture written in Go.

Requirements
------------

- [Go](http://www.golang.org/)
- A Web browser that supports WebSockets Note: Firefox is the browser
  automatically launched by the simulator.

Build and Test
-------------

Go can be complicated to build the first time. The general methodology is to
obtain the version of Go recommended for your platform, for more details see:
[http://golang.org/doc/install](http://golang.org/doc/install). Don't forget to
set a proper `$GOPATH`.

After obtaining Go, build the project on Linux by:
- `cd` to the armsim/ folder in the project.
- `export GOPATH=$GOPATH:$(pwd)` to add the armsim/ folder to your GOPATH.
- `go build -o install/armsim src/armsim.go` to build the armsim executable.
- `cd ./install/` and run `./armsim`

Note: If building fails, you many need to install the websockets package with:
`go get code.google.com/p/go.net/websocket`.

To test: `go test armsim`
  - there will be no output unless the tests fail
  - be patient, it can take 10 or 20 seconds for the test suite to run
  - (requires your GOPATH to be set)

It is entirely possible to build on Windows or Mac OS X, but you're on your own
for that.

Configuration
-------------

Command line options:
- --load: a executable file to load (if running in command line mode, this file
  will also be automatically evaluated and run)
- --log: a file name defining the location of log file (default: STDERR)
- --gui: (boolean) whether to launch the GUI or not (default: true)
- --mem: (integer) size of the memory for the simulator in bytes
- --trace: (boolean) whether or not to output a trace file (always trace.log)

You can also use `2>` to redirect most of the log output, as well.

User Guide
---------

To run the project, `cd` to the install/ directory of the project and
run the armsim executable.

When run in command line mode (`--gui=false`), the simulator will load the file
specified by `--load` and attempt to use it as an ELF structured executable. If
the loading is successful, the simulator will being simulating an execution of
the ARM instructions within the file, providing output as it goes.

When run in GUI mode (default), the simulator fires up a web server on port 4567
and attempts to run `firefox http://localhost:4567`. If this is unsuccessful,
manually navigating to [http://localhost:4567/](http://localhost:4567/) should
bring up the GUI. In the GUI you will see an array of panels and buttons (we
all like buttons!).

- Buttons
  - Load: opens a prompt for a filename and asks the simulator to open that file
    - Note: This is important to understand. The filename *must be* relative to where the
    executable is! If you downloaded the complete project: `../test/test1.exe`
    should be a valid executable.
  - Start: begins execution of the loaded file, updates the panels after execution
    has finished
  - Step: executes one step of the program, updates the panels
  - Stop/Break: ends execution of the program midstream (hey, maybe those 1,000,000
    instructions were just a few too many!), updates the panels
  - Reset: reloads the file and starts over
  - Tracing On/Off: turns the trace.log file on and off
- Panels
  - Instructions: shows the instructions that are close to the current instruction
    (currently shows dummy data)
  - Memory: shows the full contents of memory, you can even search for a specific
    address
  - Terminal: (not implemented) will eventually show output and allow input to
    the programs on the simulator
  - Flags: shows the status of the four CPSR flags (hint: if there's nothing there
    they aren't active; if there's a flag, it's true.)
  - Registers: shows the contents of r0 - r14, plus r15, the program counter
  - Stack: (not implemented) will show the top few spots in the stack

Bugs
----

- Any errors encountered by the simulator on a load are not explicitly passed to the GUI;
  this will be fixed in a future version.
- The GUI state could become corrupt if the right set of circumstances occur.
- Keyboard shortcuts enable certain commands to be run even when they should be
  disabled. This will be fixed in a future version.
