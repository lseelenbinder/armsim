typedef unsigned char uchar;
typedef unsigned int uint;

int main()
{
    int i = 0;
      
    puts("Enter a starting number:");
    i = readint();
    while (i > 0) {
      writeint(i);
      puts("...\n");
      i--;
    }
  
    asm("swi 0x11");
}