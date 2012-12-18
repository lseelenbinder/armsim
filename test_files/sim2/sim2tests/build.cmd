copy ..\sim2os\armos*.o

call dogcc cmp
call dogcc branch

arm-none-eabi-gcc.exe -c start.c -o start.o  -nostdlib -fno-builtin -nostartfiles -nodefaultlibs  -mcpu=arm7tdmi


call makec sieve 
call makec_no_io sieve
call makec_no_io locals

call makec sieve
call makec pointers 

call makec mersenne 
call makec_no_io mersenne
call makec quicksort 
call makec_no_io quicksort

:done