

arm-none-eabi-gcc.exe -c %1.c -o %1.o  -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi

arm-none-eabi-ld -T linker_no_os.ld -n -e main -o %1_no_io.exe %1.o 
arm-none-eabi-objdump -d %1_no_io.exe > %1_no_io.lst
