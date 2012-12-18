@@@ Demonstrates pc-relative addressing

.global _start
.text
_start:
  mov r0, pc
  ldr r0, [pc]
  ldr r0, [pc, #-8]
  swi 0x11
  