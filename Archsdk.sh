#!/bin/bash

ROOT="dist/tribc"

echo "ROOT is :" $ROOT
echo "Version :" $1

mkdir -p $ROOT 
mkdir -p $ROOT/core
mkdir -p $ROOT/lib
mkdir -p $ROOT/doc
mkdir -p $ROOT/test
mkdir -p $ROOT/inc

cp inc/inc.go $ROOT/inc

cp core/acc.go $ROOT/core
cp core/shield.go $ROOT/core

cp lib/encry.go $ROOT/lib
cp lib/utils.go $ROOT/lib

cp doc/*.md $ROOT/doc

cp test/test_acc.go $ROOT/test
cp test/test_perf.go $ROOT/test
cp test/rpclient.py $ROOT/test

cp trias_accs.go $ROOT
cp README_ACC.md $ROOT/README.md

tar cfz trias_acc_$1.tar.gz $ROOT
