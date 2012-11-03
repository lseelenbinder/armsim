@@@ Test cmp and b

.equ SWI_Exit, 0x11 @ Stop execution
.global _start
.text
_start:
  mov r0, #0x5
  cmp r0, #0x3  ; C = 1 
  cmp r0, #0xff000000 ; C = 0 ; 
  cmp r0, #0x7f000000 ; C = 0
  cmp r0, #0x7 ; C = 0
  mov r0, #0xfe000000 
  cmp r0, #0x3 ; C = 1
  cmp r0, #0xff000000 ; C = 0

  mov r0, #0x5
  mov r1, #-3
  ldr r2, =-2147483647
  mov r3, #-4
  cmp r0, #0x3 ; V = 0
  cmp r0, #0x7 ; V = 0
  cmp r0, r1 ; V = 0
  cmp r0, r2 ; V = 1
  cmp r2, r0 ; V = 1
  cmp r1, r0 ; V = 0
  cmp r1, r3 ; V = 0
  cmp r3, r1 ; V = 0

  cmp r3, r3 ; Z = 1

swi SWI_Exit @ stop executing: end of program
