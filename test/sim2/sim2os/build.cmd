
arm-none-eabi-gcc.exe -c armos_asm.s -o armos_asm.o  -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi
arm-none-eabi-gcc.exe -c armos.c -o armos.o  -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi
arm-none-eabi-gcc.exe -c prog.c -o prog.o  -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi


arm-none-eabi-ld -T linker_separate_os.ld  -n -e main -o prog.exe armos_asm.o armos.o prog.o 
arm-none-eabi-objdump -d prog.exe > prog.lst
rem arm-none-eabi-ld -T linker.ld  -n -e _start -o prog.exe armos_asm.o armos.o prog.o 

