@@@ C level instruction tests

.global _start
.text
_start:
  mov r0, #724
  mov r1, #0xa1000000
  mov r2, r0
  mov r2, pc
  mov r2, r1, asr #2
  mov r2, r1, lsr #2
  mov r2, r1, lsl #1
  mov r2, r0, ror #4
  mvn r3, #1
  mov r4, #4
  add r5, r4, #-3
  sub r5, r4, #3
  rsb r5, r4, #3
  and r2, r0, #0xff
  orr r2, r0, #0x12
  eor r2, r0, #732
  bic r2, r0, #0xffffff00
  mov r2, #2
  mul r5, r1, r2
  swi 0x11
  