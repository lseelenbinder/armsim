SECTIONS
{
     
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