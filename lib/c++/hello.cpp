//#include <iostream>
#include <stdio.h>
//#include <dlfcn.h>

#include "hello.h"
#include "../inc/zksnark.hpp"
//using namespace std;



void myhello(char *str) {
   printf("%s\n", str);
}

int add(int a,int b){
	return a+b;
}

void test(){
    testProv();
}

/*

void r_testProv()
{
    void* handle;
    typedef int (*FPTR)();

    handle = dlopen("/home/lab8/go/src/trias/lib/libzksnark.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "testProv");
    (*fptr)();
}
*/


/*

int r_add(int a,int b)
{
    void* handle;
    typedef int (*FPTR)(int,int);

    handle = dlopen("/home/lab8/go/src/trias/lib/libhello.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "add");

    int result = (*fptr)(a,b);

    return result;
}

void r_myhello(char *str)
{
    void* handle;
    typedef int (*FPTR)(char *str);

    handle = dlopen("/home/lab8/go/src/trias/lib/libhello.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "myhello");

    (*fptr)(str);
}
*/



