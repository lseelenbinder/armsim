@echo off
setlocal
path=c:\programs\unixutils;%path%
arm-none-eabi-objdump -d %1.exe > %1.lst
@echo Instructions used in %1.exe:
egrep "^[ ]+[0-9a-f]+:" %1.lst | cut -f 3  | sort -u
