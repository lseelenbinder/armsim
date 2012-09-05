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

## Loader (Total time: 0.5 hours)
1. Setup package structure (Time: 0.5 hours)
  - Go has a very smart package system. Unfortunately, due to the relative
    newness of the language, implementing under the Go "standard" was difficult
    to accomplish. After much googling and reading, I figured out how to order
    my project according to recommended standards. The key was temporarily
    setting `$GOPATH`.
  - References:
    - [http://golang.org/doc/code.html]()
    - [http://golang.org/doc/effective_go.html]()
