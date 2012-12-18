@@@ Test branches

.equ SWI_Exit, 0x11 @ Stop execution
.global _start
.text
_start:
  mov r5, #0 @ counts the number of branches not taken (for trace)
  mov r0, #0
  mov r1, #1
  b label1
  add r5, r5, #1
label1:
  cmp r0, r0
  beq label2   @taken
  add r5, r5, #1
label2:
  cmp r0, r1
  beq label3   @not taken
  add r5, r5, #1
label3:
  cmp r0, r0
  bne label4  @not taken
  add r5, r5, #1
label4:
  cmp r0, r1
  bne label5  @taken
  add r5, r5, #1
label5:
  mov r0, #0x5
  cmp r0, #0x3  @ C = 1
  bcs label6  @ taken
  add r5, r5, #1
label6:
  cmp r0, #0xff000000  @ C = 0
  bcs label7 @not taken
  add r5, r5, #1
label7:
  cmp r0, #0x3  @ C = 1
  bcc label8  @not taken
  add r5, r5, #1
label8:
  cmp r0, #0xff000000  @ C = 0
  bcc label9 @taken
  add r5, r5, #1
label9:
  cmp r0, #0x7
  bmi label10 @taken
  add r5, r5, #1
label10:
  cmp r0, #0x3
  bmi label11 @not taken
  add r5, r5, #1
label11:
  cmp r0, #0x7
  bpl label12 @not taken
  add r5, r5, #1
label12:
  cmp r0, #0x3
  bpl label13 @taken
  add r5, r5, #1
label13:
  ldr r2, =-2147483647
  cmp r0, #0x3 @ V = 0
  bvs label14 @not taken
  add r5, r5, #1
label14:
  cmp r0, r2 @ V = 1
  bvs label15 @taken
  add r5, r5, #1
label15:
  cmp r0, #0x3 @ V = 0
  bvc label16 @taken
  add r5, r5, #1
label16:
  cmp r0, r2 @ V = 1
  bvc label17 @not taken
  add r5, r5, #1
label17:
  cmp r0, #0x3 @ C = 1
  bhi label18 @taken
  add r5, r5, #1
label18:
  cmp r0, #0x5 @ Z = 1
  bhi label19 @not taken
  add r5, r5, #1
label19:
  cmp r0, #0x3 @ C = 1
  bls label20 @not taken
  add r5, r5, #1
label20:
  cmp r0, #0x5 @ Z = 1
  bls label21 @taken
  add r5, r5, #1
label21:
  cmp r0, #0x3 @ V = 0, N = 0
  bge label22 @taken
  add r5, r5, #1
label22:
  cmp r0, #0x7 @ V = 0, N = 1
  bge label23 @not taken
  add r5, r5, #1
label23:
  cmp r0, #0x7 @ V = 0, N = 1
  blt label24 @taken
  add r5, r5, #1
label24:
  cmp r0, #0x3 @ V = 0, N = 0
  blt label25 @not taken
  add r5, r5, #1
label25:
  cmp r0, #0x3 @ V = 0, N = 0
  bgt label26 @taken
  add r5, r5, #1
label26:
  cmp r0, #0x5 @ Z = 1
  bgt label27 @not taken
  add r5, r5, #1
label27:
  cmp r0, #0x3 @ V = 0, N = 0
  ble label28 @not taken
  add r5, r5, #1
label28:
  cmp r0, #0x5 @ Z = 1
  ble label29 @taken
  add r5, r5, #1
    
label29:
  bl mysub
  swi SWI_Exit


mysub:
  bx lr
  
  