typedef int T;
 
void quicksort(T* data, int N)
{
  int i, j;
  T v, t;
 
  if( N <= 1 )
    return;
 
  // Partition elements
  v = data[0];
  i = 0;
  j = N;
  for(;;)
  {
    while(data[++i] < v && i < N) { }
    while(data[--j] > v) { }
    if( i >= j )
      break;
    t = data[i];
    data[i] = data[j];
    data[j] = t;
  }
  t = data[i-1];
  data[i-1] = data[0];
  data[0] = t;
  quicksort(data, i-1);
  quicksort(data+i, N-i);
}

#define NUM_ELEM(arr) (sizeof(arr) / sizeof(arr[0]))
int nums[] = { 200, -15, 30, 45, 102, -2, 19 };
    
int main()
{

    int i;
    
    quicksort(nums, NUM_ELEM(nums));
    
    #ifdef IO
    for (i = 0; i < NUM_ELEM(nums); ++i) {
        writeint(nums[i]);
    }
    #endif
    
    asm("swi 0x11");
}