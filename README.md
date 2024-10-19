# anime-petition-service
动漫请愿的后端服务

## 项目介绍
由于该平台为动漫众筹请愿平台，代币只允许购买，不允许出售。并且NFT不允许私自出售，只允许在平台方出售和购买。用户在购买Token 和 NFT 时产生的gas费由平台方承担，但平台方会抽取NFT售价的1%作为平台手续费来补偿gas费用
为了保证众筹的资金安全并在指定时间内作者完成作品后才允许提现，平台方在铸造初始代币时会设置锁仓期，锁仓期内不允许提现，只有作者完成作品后才解锁
保证平台正常运营，每次平台会在铸币数量的基础上，再铸造10%的代币作为平台方的流动资金

## 用户点击购买代币
用户点击购买代币按钮，前端将用户的地址和代币数量发送给后端，后端将代币数量添加到用户的账户中，并返回一个唯一的订单号
代币价值和 USDT 1:1

```
GET /user/buy_token?username=xxx&signature=xxx&token_amount=100
```
username: 用户钱包地址
signature: 用户签名，由用户的私钥生成的哈希值，签名原始数据为 0x19 + 0x01 + chainId + username + "buy"

```
200 OK

{"status": 200, "message": "buy token success", "data": {"token_num": "300", "token_name": "AnimePetition", "prtition_num": 30}}
```
token_num：代币数量
token_name：代币名称
prtition_num：请愿数量

```
401 Unauthorized
{"status": 401, "message": "buy token failed, signature error"}
```

## 用户发起请愿
用户发起请愿，前端将用户的地址和请愿内容发送给后端，后端将请愿内容存储到数据库中，并返回一个唯一的请愿号

```
GET /user/create_petition?username=xxx&signature=xxx&petition_id=xxx
```
username: 用户钱包地址
signature: 用户签名，由用户的私钥生成的哈希值，签名原始数据为 0x19 + 0x01 + chainId + username + "create_petition"
petition_id: 请愿的角色id

```
200 OK
Set-Cookie: anime_prtition_cookie=xxx; Expires=xxxx; Path=/; Secure; HttpOnly

{"status": 200, "message": "create petition success", "data": {"token_num": "290", "token_name": "AnimePetition", "prtition_num": 31, "NFT_id": "12"}}
```
token_num：代币数量
token_name：代币名称
prtition_num：请愿数量

```
401 Unauthorized
{"status": 401, "message": "buy token failed, signature error"}
```

```
403 Forbidden
{"status": 403, "message": "create petition failed, token not enough, need 10 tokens", "data": {"token_num": "5", "token_name": "AnimePetition", "prtition_num": 31}}
```

## NFT出售
NFT出售，前端将NFT的id和价格发送给后端，后端将NFT的价格存储到数据库中，并返回一个唯一的NFT出售号

```
GET /user/sell_nft?username=xxx&signature=xxx&NFT_id=xxx&price=1000
```
username: 用户钱包地址
signature: 用户签名，由用户的私钥生成的哈希值，签名原始数据为 0x19 + 0x01 + chainId + username + "sell_nft"
NFT_id: NFT的id
price: NFT的价格

```
200 OK

{"status": 200, "message": "Successfully registered NFT for sale", "data": {"sell": true}}
```
sell: NFT已经被注册为出售状态

```
401 Unauthorized
{"status": 401, "message": "buy token failed, signature error"}
```

```
403 Forbidden
{"status": 403, "message": "NFT has been registered for sale", "data": {"sell": true, "NFT_onwer": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "price": "1000", "NFT_id": "12"}}
```

## 购买NFT
用户购买NFT，前端将用户的地址和NFT id 发送给后端，后端将NFT id 和用户的地址绑定

```
200 OK

{"status": 200, "message": "Successfully registered NFT for sale", "data": {"sell": false}}
```
sell: NFT未被注册为出售状态

```
401 Unauthorized
{"status": 401, "message": "buy NFT failed, signature error"}
```

```
403 Forbidden
{"status": 403, "message": token not enough, need 1000 tokens", "data": {"sell": true, "NFT_onwer": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "price": "1000", "NFT_id": "12"}}  
```

