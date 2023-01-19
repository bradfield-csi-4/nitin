#include "csapp.h"

volatile long cnt = 0;
pthread_mutex_t m1;
pthread_mutex_t m2;

void *thread1(void *vargp);
void *thread2(void *vargp);

int main(void) {
    pthread_t tid1, tid2;
    pthread_mutex_init(&m1, NULL);
    pthread_mutex_init(&m2, NULL);

    Pthread_create(&tid1, NULL, thread1, NULL);
    Pthread_create(&tid2, NULL, thread2, NULL);

    Pthread_join(tid1, NULL);
    Pthread_join(tid2, NULL);

    printf("NO DEADLOCK!");
}

void *thread1(void *vargp) {
    for (int i = 0; i < 10000; i++) {
        pthread_mutex_lock(&m1);
        pthread_mutex_lock(&m2);
        printf("thread1: %ld\n", ++cnt);
        pthread_mutex_unlock(&m2);
        pthread_mutex_unlock(&m1);
    }

    return NULL;
}

void *thread2(void *vargp) {
    for (int i = 0; i < 10000; i++) {
        pthread_mutex_lock(&m2);
        pthread_mutex_lock(&m1);
        printf("   thread2: %ld\n", ++cnt);
        pthread_mutex_unlock(&m1);
        pthread_mutex_unlock(&m2);
    }
    return NULL;
}