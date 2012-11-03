setlocal
path=c:\programs\unixutils;%path%
arm-none-eabi-objdump -d %1.exe > %1.lst
grep "^ ...:" %1.lst | cut -f 3  | sort -u
