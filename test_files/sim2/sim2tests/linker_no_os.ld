SECTIONS
{
    . = 0x001000;

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