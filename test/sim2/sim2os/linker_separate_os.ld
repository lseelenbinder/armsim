SECTIONS
{
    . = 0x00000;

    .ostext ALIGN (0x04) :
    {
        armos*.o(.text)
        armos*.o(.rodata)
        armos*.o(.data)
    } 

    .osbss :
    {
        sbss = .;
        armos*.o(COMMON)
        armos*.o(.bss)
        ebss = .;
    }
    
     . = 0x002000;
     
    .text ALIGN (0x04) :
    {
        *(.text)
        *(.rodata)
        *(.data)        
    }

    .bss :
    {
        sbss = .;
        *(COMMON)
        *(.bss)
        ebss = .;
    }

}