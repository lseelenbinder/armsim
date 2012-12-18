typedef unsigned char uchar;
typedef unsigned int uint;
typedef unsigned int size_t;


// Execution begins here
void start()
{
  main();
  asm("swi 0x11");
}


// C library

// output <c> to stdout
void putc(uchar c)
{
    asm("" : : : "r5"); // force use of stm for push
    // by arm calling conventions, value of c is in r0
    asm("swi 0");
}

// Read a line from stdin, store it in <buf>, null-terminated, with \r.
// Read up to <maxlen> chars. 
void getline(char *buf, size_t maxlen)
{
    asm("" : : : "r5"); // force use of stm for push
    asm("mov r2, r1");
    asm("mov r1, r0");
    asm("swi 0x6a");
}

// output null-terminated <str> to stdout
void puts(uchar *str)
{
    while (*str) {
        putc(*str++);
    }
}

// reads an integer from keyboard and returns it
int readint()
{
  uchar buf[40];
  uchar *ptr = buf;
  int result = 0;
  int i;
  int sign = 1;
  
  getline(buf, sizeof(buf));
  if (*ptr == '-') {
    ++ptr;
    sign = -1;
  }
  
  while (*ptr && *ptr != '\r') {
    uchar c = *ptr;
    result = result * 10 + (c - '0');
    ++ptr;
  }
  result = result * sign;
  return result;
  
}

// From http://discuss.joelonsoftware.com/default.asp?interview.11.401113.6
// returns <num> / <denom>
int divide(int num, int denom)
{
  asm("" : : : "r5"); // force use of stm for push

  int a=0, b=0;
  int i= 31; // CAREFUL: works only on int=32-bit machine!
  /* Work from leftmost to rightmost bit in numerator */
  while(i>=0) {
    /* appends one bit from numerator to a */
    a = (a << 1) + ((num & (1 << i)) >> i);
    b = b << 1;
    //printf("After shifting a=%d and b=%d\n",a,b);
    if (a >= denom) {
      a -= denom;
      b++;
    }
    //printf("After subtraction a=%d and b=%d\n",a,b);
  i--;
  }
  return b;
}


// converts <num> to base 10, stores in <result>
void itoa(int num, uchar *result) {
  uchar buf[20] = "0\n";
  uchar *pos = buf;
  uchar *writeptr = result;
  int numWritten;
 
  // Handle negative numbers
  if (num < 0) {
    *writeptr++ = '-';
    num = -num;
  }
  
  if (num > 0) {
      
    // Build the number in reverse order
    while (num > 0) {
      uint quotient, remainder;
      
      quotient = divide(num, 10);
      remainder = num - (quotient * 10);
      *pos++ = remainder + '0';
      num = quotient;
    }
    pos--;
    
    // Now we need to copy the results into the output buffer, reversed
    while (pos > buf) {
      *writeptr++ = *pos--;
    }
    *writeptr++ = *pos;
    *writeptr++ = 0;
  } else {
    // number is 0
    *writeptr++ = '0';
    *writeptr++ = '\0';
  }
  
}

// output <num> in base 10
void writeint(int num) {
  uchar buf[20];
  itoa(num, buf);
  
  puts(buf);
  
}

// ------------------------------------------------------
// INTERNAL OS ROUTINES
// ------------------------------------------------------

// keyboard interrupt handler 

uint kbuf_count = 0;
uint kbuf_queue_start = 0;
uchar kbuf[10];
int kbd_enabled = 1;

void internal_putc(char c) {
  asm("" : : : "r5"); // force use of stm for push

  asm("mov r1, #0x100000");
  asm("strb r0, [r1]");
}

void kbdinthandler()
{
    uchar *kbdevice = (uchar *)0x100001;
    uchar *kbuf_c;
    
    if (!kbd_enabled || kbuf_count >= sizeof(kbuf))
        return;
    
    uint kbuf_index = kbuf_queue_start + kbuf_count;
    if (kbuf_index >= sizeof(kbuf)) 
        kbuf_index = kbuf_index - sizeof(kbuf);
    kbuf[kbuf_index] = *kbdevice;
    
    internal_putc(kbuf[kbuf_index]);
    
    ++kbuf_count;
    
}

uchar getchar()
{
    asm("" : : : "r5"); // force use of stm for push

    uchar c;
    
    while (kbuf_count == 0)  
        ;
    
    kbd_enabled = 0;
    c =  kbuf[kbuf_queue_start];
    ++kbuf_queue_start;
    if (kbuf_queue_start >= sizeof(kbuf))
        kbuf_queue_start = 0;
    --kbuf_count;
    kbd_enabled = 1;
    return c;
}

// implements SWI 0x6a
void swi_getline(uchar *buf, uint maxlen)
{
    uint i = 0;
    uchar c = 0;
    while (i < maxlen-1 && c != '\r') {
        
        buf[i] = c = getchar();        
        ++i;
    }
    buf[i] = '\0';
}

