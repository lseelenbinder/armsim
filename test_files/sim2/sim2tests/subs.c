// pointers.c

int add(int a, int b)
{
  int result = a + b;
  return result;
}

int end  = 5;

int main()
{
  int num = 0;
  int result; 
  while (num < end) {
    result = add(num, num+1);
  }
  
  asm("swi 0x11\n"); // exit
}