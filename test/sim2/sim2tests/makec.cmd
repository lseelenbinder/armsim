

arm-none-eabi-gcc.exe -DIO -c %1.c -o %1.o  -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi

arm-none-eabi-ld -T linker_separate_os.ld -n -e start -o %1.exe armos_asm.o armos.o %1.o 
arm-none-eabi-objdump -d %1.exe > %1.lst
