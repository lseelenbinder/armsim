// locals.c

int sub(int z, int y, int x, int n, int o, int q) {
      asm("" : : : "r5"); // force use of stm for push

  int a=1, b=2, c=3, d=4, e=5, f=6, g=7, h=8, i=9, j=10, k=11, l=12, m=13;
  
  int s = a + b + c + d + e + f + g + h + i + j + k + l + m + n + o + q + x + y + z;
  
  return s;
}

int main()
{
  int a = sub(10, 20, 30, 40, 50, 60);
  
  asm("swi 0x11");
}