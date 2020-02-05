# go-publicChain-study
go语言版区块链学习笔记````

demo:https://andersbrownworth.com/blockchain/

文档: https://jeiwan.net/


用到的插件:
go get "github.com/boltdb/bolt"  数据库
go get "golang.org/x/crypto/ripemd160"  加密包


创建一个钱包地址:
1. 生成一对公钥和私钥
2. 想获取地址，可以通过公钥进行base58编码
3. 想要别人给我转账，把地址给别人，别人将地址进行反编码变成公钥，将公钥和数据进行签名
4. 通过私钥进行解密，解密是单方向的，只有拥有私钥的人才能进行解密

