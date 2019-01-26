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

cp doc/*.* $ROOT/doc

cp test/test_acc.go $ROOT/test
cp test/test_perf.go $ROOT/test
cp test/rpclient.py $ROOT/test

cp -R wasm $ROOT

cp trias_accs.go $ROOT
cp README_ACC.md $ROOT/README.md
cp buildwasm.sh $ROOT
cp http.go $ROOT
cp triacc_wasm.go $ROOT
cp triacc_wasm.wasm $ROOT
cp wasm_exec.html $ROOT
cp wasm_exec.js $ROOT

tar cfz trias_acc_$1.tar.gz --exclude=.git $ROOT
