	.cpu arm7tdmi
	.fpu softvfp
	.eabi_attribute 20, 1
	.eabi_attribute 21, 1
	.eabi_attribute 23, 3
	.eabi_attribute 24, 1
	.eabi_attribute 25, 1
	.eabi_attribute 26, 1
	.eabi_attribute 30, 6
	.eabi_attribute 34, 0
	.eabi_attribute 18, 4
	.file	"sieve.c"
	.comm	sieve,6252,4
	.text
	.align	2
	.global	main
	.type	main, %function
main:
	@ Function supports interworking.
	@ args = 0, pretend = 0, frame = 8
	@ frame_needed = 1, uses_anonymous_args = 0
	stmfd	sp!, {fp, lr}
	add	fp, sp, #4
	sub	sp, sp, #8
	bl	make
	mov	r3, #0
	str	r3, [fp, #-12]
	mov	r3, #0
	str	r3, [fp, #-8]
	b	.L2
.L4:
	ldr	r0, [fp, #-8]
	bl	isprime
	mov	r3, r0
	cmp	r3, #0
	beq	.L3
	ldr	r3, [fp, #-12]
	add	r3, r3, #1
	str	r3, [fp, #-12]
.L3:
	ldr	r3, [fp, #-8]
	add	r3, r3, #1
	str	r3, [fp, #-8]
.L2:
	ldr	r2, [fp, #-8]
	ldr	r3, .L5
	cmp	r2, r3
	ble	.L4
	ldr	r3, [fp, #-12]
	mov	r0, r3
	sub	sp, fp, #4
	ldmfd	sp!, {fp, lr}
	bx	lr
.L6:
	.align	2
.L5:
	.word	100000
	.size	main, .-main
	.align	2
	.global	make
	.type	make, %function
make:
	@ Function supports interworking.
	@ args = 0, pretend = 0, frame = 16
	@ frame_needed = 1, uses_anonymous_args = 0
	@ link register save eliminated.
	str	fp, [sp, #-4]!
	add	fp, sp, #0
	sub	sp, sp, #20
	mov	r3, #1
	str	r3, [fp, #-12]
	b	.L8
.L12:
	ldr	r3, [fp, #-12]
	mov	r2, r3, lsr #5
	ldr	r3, .L13
	ldr	r2, [r3, r2, asl #2]
	ldr	r3, [fp, #-12]
	and	r3, r3, #31
	mov	r3, r2, lsr r3
	and	r3, r3, #1
	cmp	r3, #0
	bne	.L9
	ldr	r3, [fp, #-12]
	mov	r3, r3, asl #1
	add	r3, r3, #1
	str	r3, [fp, #-16]
	ldr	r3, [fp, #-12]
	add	r3, r3, #1
	ldr	r2, [fp, #-12]
	mul	r3, r2, r3
	mov	r3, r3, asl #1
	str	r3, [fp, #-8]
	b	.L10
.L11:
	ldr	r3, [fp, #-8]
	mov	r2, r3, lsr #5
	ldr	r3, [fp, #-8]
	mov	r1, r3, lsr #5
	ldr	r3, .L13
	ldr	r1, [r3, r1, asl #2]
	ldr	r3, [fp, #-8]
	and	r3, r3, #31
	mov	r0, #1
	mov	r3, r0, asl r3
	orr	r1, r1, r3
	ldr	r3, .L13
	str	r1, [r3, r2, asl #2]
	ldr	r2, [fp, #-8]
	ldr	r3, [fp, #-16]
	add	r3, r2, r3
	str	r3, [fp, #-8]
.L10:
	ldr	r2, [fp, #-8]
	ldr	r3, .L13+4
	cmp	r2, r3
	bls	.L11
.L9:
	ldr	r3, [fp, #-12]
	add	r3, r3, #1
	str	r3, [fp, #-12]
.L8:
	ldr	r3, [fp, #-12]
	cmp	r3, #159
	bls	.L12
	add	sp, fp, #0
	ldmfd	sp!, {fp}
	bx	lr
.L14:
	.align	2
.L13:
	.word	sieve
	.word	49999
	.size	make, .-make
	.align	2
	.global	isprime
	.type	isprime, %function
isprime:
	@ Function supports interworking.
	@ args = 0, pretend = 0, frame = 8
	@ frame_needed = 1, uses_anonymous_args = 0
	@ link register save eliminated.
	str	fp, [sp, #-4]!
	add	fp, sp, #0
	sub	sp, sp, #12
	str	r0, [fp, #-8]
	ldr	r3, [fp, #-8]
	cmp	r3, #2
	beq	.L16
	ldr	r3, [fp, #-8]
	cmp	r3, #2
	ble	.L17
	ldr	r3, [fp, #-8]
	and	r3, r3, #1
	and	r3, r3, #255
	cmp	r3, #0
	beq	.L17
	ldr	r3, [fp, #-8]
	sub	r3, r3, #1
	mov	r2, r3, asr #6
	ldr	r3, .L19
	ldr	r2, [r3, r2, asl #2]
	ldr	r3, [fp, #-8]
	sub	r3, r3, #1
	mov	r3, r3, asr #1
	and	r3, r3, #31
	mov	r3, r2, lsr r3
	and	r3, r3, #1
	cmp	r3, #0
	bne	.L17
.L16:
	mov	r3, #1
	b	.L18
.L17:
	mov	r3, #0
.L18:
	mov	r0, r3
	add	sp, fp, #0
	ldmfd	sp!, {fp}
	bx	lr
.L20:
	.align	2
.L19:
	.word	sieve
	.size	isprime, .-isprime
	.ident	"GCC: (Sourcery CodeBench Lite 2012.03-56) 4.6.3"
