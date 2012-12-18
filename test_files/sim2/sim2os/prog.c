typedef unsigned char uchar;
typedef unsigned int uint;

int main()
{
    uchar buf[10];
    uint i = 0;
      
    puts("Enter your name:");
    getline(buf, 10);
    puts("Hello, ");
    puts(buf);
  
    asm("swi 0x11");
}