@@@ Demonstrates stm and ldm

.global _start
.text
_start:
  mov r1, #1
  mov r2, #2
  mov r4, #4
  stmfd sp!, {r1, r2, r4}
  mov r1, #10
  mov r3, #30
  stmfd sp!, {r1, r3}
  mov r1, #0
  mov r3, #0
  ldmfd sp!, {r1, r3}
  ldmfd sp, {r1, r2, r4}
  swi 0x11
  