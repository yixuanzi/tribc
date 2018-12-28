package main

/*
#cgo linux CFLAGS: -I../inc
#cgo linux LDFLAGS: -L../lib -lhello


#include <stdlib.h>
#include <stdio.h>
#include "hello.h"
//#include "zksnark.hpp"

int madd(int a,int b){
	return a+b;
}
*/
import "C"
import (
	"fmt"
)

func main() {

	//C.r_myhello(C.CString("Hello world C"))
	//fmt.Println(C.r_add(8,9))
	//C.testProv()
	C.myhello(C.CString("Hello world C"))
	fmt.Println(C.add(8,9))
	C.test()
	fmt.Println(C.rand())
	fmt.Println("Hello Go")
	fmt.Println(C.madd(8,9))
}


//go -> 路由动态库so，使用extern "C"强制导出go能够直接调用的函数(C) -> 功能实现库，基于zksnark库实现关于trias零知识证明金额加密功能函数(C++) -> zksnark.so (开源项目 C++）
//gcc  -o hello.o hello.cpp -c -g -Wall
//gcc -shared -o libhello.so hello.o -g -Wall -lzksnark
//try using -rpath or -rpath-link 指定加载动态库路径 （-Wl,-rpath $LIB_PATH）
//当前测试demo中关于动态链接库的加载都使用系统库路径加载
//#cgo linux LDFLAGS: -L../lib -Wl,-rpath /home/lab8/go/src/tribc/lib -lhello
