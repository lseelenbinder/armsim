# ARMSIM Project Journal
  by Luke Seelenbinder

## Total Time: 51.25

## Preparation (Total time: 2 hours)
1. Investigate GO as a possibility for the project. (Time: 2 hours)

## Bootstrapping Code (Total time: 0.75 hours)
1. Implement main() and command line flag parsing. (Time: 0.75 hours)
  - This was easily accomplished via the built in Go tools (see references). I
    spent most of my time learning Go conventions and fighting with the little
    things that vary from language to language.
  - References:
    - [http://golang.org/pkg/log/](http://golang.org/pkg/log/)
    - [http://golang.org/pkg/flag/](http://golang.org/pkg/flag/)

## Loader (Total time: 11.5 hours)
1. Setup package structure (Time: 0.5 hours)
  - Go has a very smart package system. Unfortunately, due to the relative
    newness of the language, implementing under the Go "standard" was difficult
    to accomplish. After much googling and reading, I figured out how to order
    my project according to recommended standards. The key was temporarily
    setting `$GOPATH` to include my source directory.
  - References:
    - [http://golang.org/doc/code.html](http://golang.org/doc/code.html)
    - [http://golang.org/doc/effective_go.html](http://golang.org/doc/effective_go.html)
    - [http://lmgtfy.com/?q=golang+gopath](http://lmgtfy.com/?q=golang+gopath)
2. Implement RAM (Time: 7.0 hours)
  - Notes:
    - The implementation of RAM was fairly straight-forward until I got to the
      HalfWords and Words. Testing also posed a bit of a challenge because Go is
      very type safe and intrinsically prevents many standard type-related bugs;
      hence, my test cases can be much shorter and concise.
    - Implementing HalfWords and Words proved difficult. Because the base unit is
      the byte, I needed to split up the words so I could store them
      contiguously as bytes. This required shifts, casting, and additions.
      Eventually, this was implemented with a helper function that worked for
      any number of bytes, a "multi-byte" reader and writer.
    - Implementing the TestFlag, SetFlag, and ExtractBits methods was a very
      good exercise in bitwise operations and pushed me to further learning
      binary mathematics, testing and Go.
    - Some of my time was spent learning Go best practices and testing
      procedures. After implementing RAM, I hope that the learning curve for Go
      will be surmounted.
  - References:
    - [http://golang.org/pkg/testing](http://golang.org/pkg/testing)
3. Implement ELF Loader (2.0 hours)
  - Notes:
    - After implementation of well-tested RAM, the loader was rather simple.
      I wrote the acceptance tests and without issue wrote the rest of the
      loader code. The biggest difficulty was fully understanding the ELF
      format and how to properly read binary into Go's provided structs.
  - References:
    - [http://golang.org/pkg/debug/elf](http://golang.org/pkg/debug/elf)
    - [http://golang.org/pkg/binary/encoding](http://golang.org/pkg/binary/encoding)
    - [http://golang.org/pkg/os/](http://golang.org/pkg/os/)
    - [http://golang.org/src/pkg/debug/elf/file.go](http://golang.org/src/pkg/debug/elf/file.go)
4. General Refactoring (Time: 2 hours)
  - After a few days and many lines of Go code, I saw a few things that needed
    to be changed, primarily my error handling code. I spent a few hours
    getting this cleaned-up and doing some general refactoring work.
- Notes:
  - (9/7/12) At this checkpoint, I feel like the various code and test suites
    "work". However, the code itself has much maturing to do. This is a result
    of 1) learning a new language as I go and 2) figuring out requirements as I
    go. I don't believe my code is sufficient commented, DRYed, or tested.
    However, due to the checkpoint nature of the project, I plan to improve the
    above mentioned aspects dramatically in the next few weeks.
  - (9/7/12) I have very much enjoyed learning Go. The language is very strict,
    but I've found it meldable and usable (a rare combination). The
    built-in libraries are quite sufficient; unfortunately, the number of
    applicable articles and packages are quite limited.

## Prototype (Total time: 20 hours)
1. Design Prototype (Time: 3 hours)
  - Notes:
    - Since I am using a browser-based interface, I will be using JS, HTML, and
      CSS to power the interface. I decided to use Twitter's excellent
      Bootstrap. Bootstrap will provide basic styling and layout tools for the
      graphical interface.
    - For icons, I am using Font Awesome, a free, open-source collection of SVG
      icons.
    - Using the basic design elements from Visual Studio and KDBG, I designed an
      interface. The actually coding and layout did not take very long. Design
      decisions will need to be tweaked as I connect the backend with the new
      interface.
  - References:
    - [Bootstrap](http://twitter.github.com/bootstrap/index.html)
    - [Font Awesome](http://fortawesome.github.com/Font-Awesome/)
2. Develop CPU (Time: 2 hours)
  - Notes:
    - I made a design decision to make all parts of the simulator part of one Go
      package. This simplified `import` statements, testing and code scope.
    - I also had to rename RAM to Memory at this time.
    - The CPU was straight-forward to implement. Testing was a little bit
      tricky, but not anything too bad. I included a section of constants to
      hold the positions of registers in a Memory unit. I learned a bit more
      about Go memory modeling after fighting with a nil pointer bug for 30
      minutes. (This was due to not using `new()`.)
  - References:
    - The Go docs (from here on out, you can assume I used these heavily)
3. Develop Computer (Time: 2 hours)
  - Notes:
    - I encountered no difficulties in designing the Computer class. The
      computer was the simplest section to code and test. I had to fake testing
      for the `Run()` and `Step()` methods, because they are mostly stubs at
      this stage.
4. Implement Web Server (Time: 5 hours)
  - Because I am using a web interface for this application, I need to develop
    a web server that supported both static assets and WebSockets communication.
    This was a significant undertaking. I felt like I was throughly behind and
    needed to catch up on knowledge of many ideas as I went. Thankfully, Go provides
    some basics for web servers and the Internet yielded a lot of help.
  - The main components of the web server are the static assets server and the
    WebSocket server. The various concurrency issues required me to also learn
    how Go handles concurrency, which is very unusual.
  - Note: very few of the server methods are formally commented, this is because
    1) I was running out of time and 2) they are very self-evident in purpose.
  - References:
    - [WebSocket Example](http://gary.beagledreams.com/page/go-websocket-chat.html)
    - [Go WebSocket Implementation](https://code.google.com/p/go/source/browse/?repo=net#hg%2Fwebsocket)
    - [ChessBuddy](https://github.com/tux21b/ChessBuddy) (for reference and ideas)
5. Connect Design and Prototype (7 hours)
  - I spent a very long time working on this portion of the development. I used
    a combination of jQuery and my previously developed server to connect the
    prototype simulator to the prototype GUI. The exercise was grueling; I had to
    handle everything from JSON encoding to concurrency, and there were a lot of
    design decisions to be made. I am not quite satisfied with pushing the entirety
    of RAM for every update, but I think the local server/socket configuration
    will mitigate speed issues. I am also concerned about browser memory issues,
    but modern computers should perform well enough.
  - References:
    - [Detecting Non-Printable Characters in Javascript](http://stackoverflow.com/questions/1677644/detect-non-printable-characters-in-javascript)
    - [Automating Scrolling in JS](http://flesler.blogspot.com/2007/10/jqueryscrollto.html)
      & [http://demos.flesler.com/jquery/scrollTo/](http://demos.flesler.com/jquery/scrollTo/)
    - [Keyboard Shortcuts](http://www.stepanreznikov.com/js-shortcuts/)
    - [Bootstrap Docs](http://twitter.github.com/bootstrap/)
6. Documentation (1 hour)
- Notes:
  - (9/18/12) After 12 hours straight of working on this prototype, I think I can
    safely say I don't want to see it for a while. However, all the requirements
    are met with more or less quality. I know there are little bugs to be found
    and fixed and probably big bugs, too. Because of the scope of the system, I
    don't yet feel like I have complete control over it. But, for now, it works
    unless crazy things are thrown at it, and that is what matters.

## Simulator I (Total Time: 17 hours)
1. Build Instruction Decode and Execution Pathway (Time: 6 hours)
  - Notes:
    - I needed to build a system by which I could take arbitrary instructions and
      narrow them down quickly and effciently. I ended up building a solution 
      based on a BaseInstruction class and general classes for each type of 
      instruction. By also including a BarrelShifter class, I made handling of
      operand2s very simple.
    - I found Go's interfaces very intuitive and good for the system I ended up 
      building.
    - Building the BarrelShifter class was difficult, but not nearly as difficult
      as I expected.
    - A good OOP approach to the instructions made building, testing and programming
      very simple.
  - References:
    - [ROR algorithm](http://stackoverflow.com/questions/3476969/rotate-bits-right-operation-in-ruby)
    - I'm not entirely sure where the ASR logic came from. I think it was partially my solution
      with a bit of googling mixed in.
2. Build Instructions (Total time: 11 hours)
  - Notes:
    - This process was not very difficult (though some instructions proved difficult to get just right).
    - I found comparison of the trace's output with expected to be a very good way to test, at least initally.
    - The sheer amount of information required when programming an instruction is very laborious.
      I can't imagine those that program instruction sets with real hardware.
