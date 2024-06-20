/* gcc -o level level.c */
#include <stdlib.h>

int main(void) {
  return system("ls /home/super-bobby/.passwd");
}
