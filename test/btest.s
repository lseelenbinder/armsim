@@@ B level instruction tests

.global _start
.text
_start:
  mov r1, #0xfb0
  mov r2, #0x5000
  mov r3, #0x3000
  mov r4, #8
  str r1, [r2]
  str r1, [r2, #-4]    @ 4ffc
  str r1, [r2, -r4]     @ 4ff8
  str r1, [r2, r4]      @ 5008
  str r1, [r2, r4, asr #1]  @5004
  str r1, [r2, r4, lsl #2]  @5020
  strb r1, [r2, #12]  @ 500c
  ldr r5, [r2]
  ldr r6, [r2, #-4]    @ 4ffc
  ldr r7, [r2, -r4]     @ 4ff8
  ldr r8, [r2, r4]      @ 5008
  ldr r9, [r2, r4, asr #1]  @5004
  ldr r10, [r2, r4, lsl #2]  @5020
  ldrb r11, [r2, #12]  @ 500c
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
  