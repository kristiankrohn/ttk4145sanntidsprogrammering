
// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;
// Note the return type: void*
void* thread_1function(){
  int j;
    for(j = 0; j <= 1000000; j++){
      i++;
    }
    return NULL;
}

void* thread_2function(){
  int k;
    for(k = 0; k <= 1000000; k++){
      i--;
    }
    return NULL;
}

int main(){
    pthread_t thread_1, thread_2;
    pthread_create(&thread_1, NULL, thread_1function, NULL);
    pthread_create(&thread_2, NULL, thread_2function, NULL);
    pthread_join(thread_1, NULL);
    pthread_join(thread_2, NULL);
    printf("%d\n", i);
    return 0;

}
