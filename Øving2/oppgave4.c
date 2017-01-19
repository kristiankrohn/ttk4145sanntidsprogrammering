
// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t lock_i;

// Note the return type: void*
void* thread_1function(){
  pthread_mutex_lock(&lock_i);
  int j;
    for(j = 0; j <= 1000000; j++){
      i++;
    }
    pthread_mutex_unlock(&lock_i);
    return NULL;
}

void* thread_2function(){
  pthread_mutex_lock(&lock_i);
  int k;
    for(k = 0; k <= 1000000; k++){
      i--;
    }
    pthread_mutex_unlock(&lock_i);
    return NULL;
}

int main(){
    pthread_mutex_init(&lock_i, NULL);
    pthread_t thread_1, thread_2;
    pthread_create(&thread_1, NULL, thread_1function, NULL);
    pthread_create(&thread_2, NULL, thread_2function, NULL);
    pthread_join(thread_1, NULL);
    pthread_join(thread_2, NULL);
    pthread_mutex_destroy(&lock_i);
    printf("%d\n", i);
    return 0;

}
