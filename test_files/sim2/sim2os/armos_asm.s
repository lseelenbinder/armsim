@ -------------------------------------------
@ File: armos_asm.s
@ Desc: Defines ARM vector table and assembly routines
@       for exception handling
@ -------------------------------------------

.text
.global _start
_start:
@ -------------------------------------------
@ ARM Vector Table
@ -------------------------------------------
    b do_reset  @ RESET
    mov r0, r0  @ Undefined
    b do_swi    @ SWI
    mov r0, r0  @ Prefetch abort
    mov r0, r0  @ Data abort
    mov r0, r0  
    b do_irq    @ IRQ

.org 0x30
@ -------------------------------------------
@ do_swi
@ - exception handler for SWI call
@ -------------------------------------------
do_swi:
    push {r0-r12, lr}

    @ enable IRQ

    @ mrs r8, cpsr
    @ bic r8, r8, #0x80
    @ msr cpsr, r8
    
    @ Determine 
    ldr r8,[lr, #-4]        @ load SWI instruction into r8
    bic  r8,r8,#0xff000000  @ mask out top 8 bits
    
    @ Branch to appropriate SWI handler function
    cmp r8, #0 
    beq do_putc
    cmp r8, #0x6a
    beq do_getline
    
    b swi_return    @ unknown swi op - ignore
    
do_putc:
    mov r1, #0x100000
    strb r0, [r1]
    b swi_return
    
do_getline:
    bl swi_getline   @ See armos.c
    
swi_return:
    pop {r0-r12, lr}
    movs pc, lr      @ context switch

@ -------------------------------------------
@ do_irq
@ - exception handler for IRQ 
@ -------------------------------------------
do_irq:
    sub lr, lr, #4
    push {r0-r12, lr}
    bl kbdinthandler
    pop {r0-r12, lr}
    movs pc, lr      @ context switch
    
@ -------------------------------------------
@ do_reset
@ - exception handler for RESET
@ - for demonstration only; not used in your simulator
@ -------------------------------------------
do_reset:
    @ Set up stacks
    
    mov r2, #0xc0 | 0x12
    msr cpsr, r2    @ switch to IRQ mode with interrupts disabled
    ldr sp, =0x7ff0

    mov r2, #0xc0 | 0x13
    msr cpsr, r2    @ switch to Supervisor mode with interrupts disabled
    ldr sp, =0x78f0

    mov r2, #0xc0 | 0x1f
    msr cpsr, r2    @ switch to SYS mode with interrupts disabled
    ldr sp, =0x7000

    @ enable IRQ
    mrs r1, cpsr
    bic r1, r1, #0x80
    msr cpsr, r1
    
    @ Now, branch to startup code in armos.c
    bl start  
    
    