# ARMSIM Project Journal
  by Luke Seelenbinder

## Preparation (Total time: 2 hours)
1. Investigate GO as a possibility for the project. (Time: 2 hours)

## Bootstrapping Code (Total time: 0.75 hours)
1. Implement main() and command line flag parsing. (Time: 0.75 hours)
  - This was easily accomplished via the built in Go tools (see references).
  I spent most of my time learning Go conventions and fighting with the little
  things that vary from language to language.
  - References:
    - [http://golang.org/pkg/log/]()
    - [http://golang.org/pkg/flag/]()

## Loader (Total time: 4.5 hours)
1. Setup package structure (Time: 0.5 hours)
  - Go has a very smart package system. Unfortunately, due to the relative
    newness of the language, implementing under the Go "standard" was difficult
    to accomplish. After much googling and reading, I figured out how to order
    my project according to recommended standards. The key was temporarily
    setting `$GOPATH`.
  - References:
    - [http://golang.org/doc/code.html]()
    - [http://golang.org/doc/effective_go.html]()
2. Implement RAM (Time: 4.0 hours)
  - The implementation of RAM was fairly straight-forward until I got to the
    HalfWords and Words. Testing also posed a bit of a challenge because Go is
    very type safe and intrinsically prevents many standard type-related bugs;
    hence, my test cases can be much shorter and consise.
  - Implementing HalfWords and Words proved difficult. Because the base unit is
    the byte, I needed to split up the words so I could store them
    contiguously as bytes. This required shifts, casting, and additions.
    Eventually, this was implemented with a helper function that worked for any
    number of bytes, a "multi-byte" writer and reader.
  - Note: Again, some of my time was spent learning Go best practices and testing
    procedures. After this, I hope that the learning curve for Go will be
    surmounted.
  - References:
    - [http://golang.org/pkg/testing]()
