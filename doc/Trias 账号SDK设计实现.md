# Trias 账号SDK设计实现



## 账号结构

```go
//密钥对结构
type GKey struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

//账户结构，采用了两个秘钥对，用于后期对隐私交易的支持，在不适用隐私交易的情况中，仅仅采用第一个秘钥对进行签名验证
type Account struct {
	GkeyA *GKey 
	GkeyB *GKey
}
```



## 对外功能接口设计

1. lib.GenerateRstring（i int) //生成指定长度的随机字符串
2. core.Priv2gkey(priv_s []byte) //基于私钥的byte数组返回一个GKey指针
3. core.Pub2pubKey(pub_b []byte) //基于公钥byte数组返回一个公钥指针
4. lib.AesEncrypt(pass_b, aeskey) //基于输入数据，返回加密结果和错误标志，参数一为bytes类型的待加密数据，参数二为bytes类型的密码
5. lib.AesDecrypt(Pass_b, aeskey) //基于加密数据，返回解密数据和错误标志，参数一为加密数据，参数二为密码，若返回都为nil，则解密失败



## 对外业务接口设计

1. gkey_X, err :=  core.MakeNewKey(s string) //根据输入的字符串生成一个秘钥对
2.  acc := core.Account{gkeyA,gkeyB} //构建账户结构对象
3. core.GetAddress(pub_bytes []byte) //根据公钥的byte数组，计算用户使用的地址字符串，推荐使用账户结构中GkeyA秘钥对的公钥
4. core.Save2file(&acc,"/tmp/trias_acc.json",[]byte("passwd")) //使用密码保存账号到文件
5. acc:= core.Load4file("/tmp/trias_acc.json",[]byte("passwd")) //使用密码从文件导入账号，返回一个账号对象的指针，若为空则失败
6. stext,err := core.Sign(priv,text) //根据私钥指针，对text进行签名，实践中为效率考虑，一般对签名数据的哈希进行签名，返回签名数据字符串和错误码
7. f,err:=core.Verify(text,stext,&priv.PublicKey) //输入需要校验的数据（一般为hash），签名数据，公钥指针，实现对签名有效性的检查，返回标志和错误码
8. shieldaddr, shieldpKey := core.CreateShieldAddr(&acc) // 为目标账户账户生成临时地址，账户中只需要公钥内容，私钥无需填充，返回隐私地址，对应隐私地址中当做随机数的公钥
9. core.Verify_shield(&acc, shieldaddr, shieldpKey) //根据账户数据，隐私地址，公钥随机数，检测对应隐私地址的交易是否是对应账户的
10. priv := core.Getprivkey(&acc, shieldpKey) //根据账户数据，公钥随机数，计算返回对应隐私地址的私钥，用于对对应隐私交易的支付验证



## 功能测试样例

1. 解压文件trias_sdk.tar.gz
2. 复制dist目录中的tribc文件夹到 $GOPATH/src 
3. 执行测试文件：go run test/test_acc.go

## 性能测试样例
对当前提供的所有接口中需要频繁执行的接口，分析其操作性能，统计如下：(其性能测试脚本执行 go run test/test_perf.go)

1. The core.Verify_shield performance:  8409 次/s
2. The core.Sign performance:  6220 次/s
3. The core.Verify performance:  11124 次/s
4. The lib.AesEncrypt performance:  504 MB/S
5. The lib.AesDecrypt performance:  723 MB/S


# trias 账号SDK服务化

## RPC接口设计
1. 生成一个新账号并导出保存:AccRPC.CreateAcc
    1. path:保存路径
    2. pass:加密密码

2. 返回当前账号列表:AccRPC.GetAcclist

3. 导出账号:AccRPC.ExportAcc
    1. path:导入路径
    2. pass:解密密码

3. 导入账号:AccRPC.ImportAcc
	1. addr ：账号地址
	2. path：导出账号文件夹路径，成功时将导出账号文件到 path/addr.json 
	3. pass1：导出账号的授权验证密码，和加密密码一致
	4. pass2:新的导入账号文件加密密码，当为空时，加密密码保持不变

4. 生成签名:AccRPC.Sign
    1. addr:用户地址
    2. hash:签名hash
    3. pass:用户密码（和保存密码一致）

5. 签名验证:AccRPC.Verify
    1. pubkey:公钥
    2. hash:数据hash
    3. stext：签名数据

6. 生成隐私地址，用于发起隐私交易:AccRPC.CreateShieldAddr
    1.用户地址

7. 验证隐私地址,用于遍历确认隐私交易输出是否是对应用户的:AccRPC.Verify_shield
    1. addr:用户地址
    2. shieldaddr：隐私地址
    3. shieldpkey：隐私交易随机数

8. 隐私地址签名，用于支付受隐私交易保护的utxo:AccRPC.Shield_Sign
    1. addr:用户地址
    2. pass:用户密码
    3. hash：签名hash
    4. shieldpkey：隐私交易随机数

9. 由地址返回公钥B:AccRPC.GetACCPubkB
    1. 公钥字符串

10. 由隐私地址公钥到地址:AccRPC.Shield_Pubkey2Addr
    1. 公钥字符串

-------
tips:签名检查分两步：签名有效性检查，当前签名用户地址和当前资产地址一致性检查
## RPC测试环境
1. 192.168.1.200:9876
2. 连接测试脚本 test/rpclinet.py (python3 下测试通过)

# trias 账号模块wasm兼容
## wasm 下的接口设计
在成功加载wasm模块后，浏览器存在一个对象triacc,通过此对象可以有如下的接口调用

1. triacc.CreateAcc() //p1:加密密码 p2:回调函数，并返回账号结果
   - 创建账号，回调返回加密账户数据
2. triacc.Load4Data() //p1:加密账号数据  p2:解密密码  p3:回调函数，并返回账号结果
   - 根据输入的加密账户数据，回调返回加载的账号地址
3. triacc.GetCurrentAcc() //p1:回调函数，并返回当前账号地址
   - 回调返回当前操作账号地址
4. triacc.Sign() //p1:需签名数据  p2:回调函数，并返回签名数据结构，包含公钥和签名数据
   - 使用当前账号，对签名数据进行签名，回调返回已签名数据结构，包括公钥和签名后数据 
5. triacc.Verify() //p1:公钥 p2:签名数据 p3:签名验证数据 p4:回调函数，返回验证结果
   - 基于输入公钥和签名数据进行验证，返回签名验证结果
6. triacc.CreateShieldAddr() //p1:用户对外地址 p2:回调函数，返回隐私地址相关数据
   - 基于输入账户地址，回调返回隐私地址数据结构，包括隐私地址和匹配机密随机数
7. triacc.Verify_shield() //p1:shieldaddr p2:shieldpkey p3:回调函数，返回判断状态
   - 使用当前账号，基于输入的隐私地址，机密随机数，回调返回判断当前隐私地址是否是当前账号拥有的
8. triacc.Shield_Sign() //p1:shieldpkey p2:需要签名数据 p3:回调函数，返回签名数据结构
   - 基于当前账号，基于输入机密随机数，对隐私地址花费数据进行签名，回调返回签名数据结构

## wasm 测试和triascc功能接口测试
1. 执行buildwasm.sh triacc_wasm 编译生成triacc_wasm.wasm(要求>=go 1.11.0，默认已提供编译完成的wasm文件)
2. 执行http.go 打开server服务器
3. 使用浏览器访问 localhost:8080/wasm_exec.html，点击其中的按钮进行功能测试，打开console查看调试信息