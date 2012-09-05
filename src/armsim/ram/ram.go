package ram

import "log"

type RAM struct {
  memory []byte
}

func Init(nBytes uint32) RAM {
  log.Println("Initializing", nBytes, "bytes of RAM...")
  ram := RAM{make([]byte, nBytes)}

  return ram
}
