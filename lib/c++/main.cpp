#include<stdio.h>

#include "hello.h"

int main(void)
{
    myhello("hello world\n");
    printf("this is a test call with %d\n",add(8,8));
    test();
    return 0;
}
