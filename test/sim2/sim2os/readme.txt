README
=====

To use this OS with the UV ARM simulator, do the following:
  
* Uncomment lines 31-33 in armos_asm.s
* Link using this command:
   arm-none-eabi-ld -T linker.ld  -n -e _start -o prog.exe armos_asm.o armos.o prog.o

This will cause the UV ARM simulator to start at address 0, 
run the reset handler to setup the stacks, and call main