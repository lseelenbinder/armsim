// pointers.c

int main()
{
  int data[] = {3, 5, 15, 25, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 0 };
  int *ptr;
  int i, size;
  int sum = 0;
  
  size = sizeof(data) / sizeof(data[0]);

  ptr = data;
  while (*ptr++)
    ;

  ptr = data;
  while (*(++ptr))
    ;      

  for (i = 0; i < size; ) {
    sum += data[i++];
  }
  
  for (i = 0; i < size; ) {
    sum += data[++i];
  }
}