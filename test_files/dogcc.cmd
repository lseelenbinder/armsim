@if exist %1.s set infile=%1.s
@if exist %1.c set infile=%1.c
arm-none-eabi-gcc.exe -c %infile%   -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi

arm-none-eabi-ld %2 %3 %4 %5 %6 %7 %8 %9 -T linker.ld -n -e _start -o %1.exe %1.o