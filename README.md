# WBA-BC-Project-01

# ๐ ๋ชฉ์ฐจ 
 - [์๊ฐ](#์๊ฐ) 
 - [์ฌ์ฉ๋ฒ](#์ฌ์ฉ๋ฒ)
 - [๋๋ ํ ๋ฆฌ ๊ตฌ์กฐ](#๋๋ ํ ๋ฆฌ-๊ตฌ์กฐ)
 - [๊ฐ๋ฐ ํ๊ฒฝ](#๊ฐ๋ฐ-ํ๊ฒฝ)
 - [์ฌ์ฉ ๊ธฐ์ ](#์ฌ์ฉ-๊ธฐ์ )
 - [ERD](#erd)
 - [์๋ฒ ์ํคํ์ฒ](#์๋ฒ-์ํคํ์ฒ) 
 - [Api ๋ช์ธ์](#api-๋ช์ธ์)
    - [๋ฐ๋ชฌ ์๋ฒ](#๋ฐ๋ชฌ์๋ฒ)
        - [์ด๋๋ฆฌ์ ๋ฐ๋ชฌ](#์ด๋๋ฆฌ์)
        - [ํด๋ ์ดํผ ๋ฐ๋ชฌ](#ํด๋ ์ดํผ)
        - [์๋ฏน์ค ๋ฐ๋ชฌ](#์๋ฏน์ค)
    - [๋ฉํฐ ์ฒด์ธ ์ง๊ฐ](#๋ฉํฐ์ฒด์ธ์ง๊ฐ)
        - [๋ก๊ทธ์ธ](#๋ก๊ทธ์ธ)
        - [์ฝ์ธ์ ์ก](#์ฝ์ธ์ ์ก)
       
    
# ์๊ฐ 
 - ๊ตฌ๊ธ ์์ ๋ก๊ทธ์ธ์ ์ด์ฉํ ๋ฉํฐ ์ฒด์ธ ์ง๊ฐ
 ๊ฑฐ๋์๊ฐ ์๋ ๊ฐ์ธ์ง๊ฐ์ ์ด์ฉ ์์ ์ธ ์ ๊ท ์ ์ ๋ ์ง๊ฐ์์ฑ๋ถํฐ ์ด๋ ค์์ ๋๋๋ค. ๊ทธ๋์ ๊ตฌ๊ธ ๋ก๊ทธ์ธ์ ํ๋ฉด ์๋์ผ๋ก ์ง๊ฐ์ ์์ฑํ๊ณ , ์ ์ ๊ฐ ๋๋ชจ๋ ๋จ์ด์ ์ง๊ฐ ํค๋ฅผ ๊ด๋ฆฌํ์ง ์์๋ ๋๋ ์ง๊ฐ์ด ์์ผ๋ฉด ์ข๊ฒ ๋ค๊ณ  ์๊ฐํ์ต๋๋ค. 

 ์ ํฌ ํ๋ก์ ํธ๋ ๊ธฐ์กด ์์๊ณผ ๋์ผํ๊ฒ ๋ฉ์ผ์ฃผ์์ ํจ์ค์๋๋ง ์์ผ๋ฉด ๋๊ตฌ๋ ์ฌ์ฉ ๊ฐ๋ฅํ ์ ์  ์นํ์ ์ธ ์ง๊ฐ์ ์ ๊ณตํฉ๋๋ค.

## TEAM Mallet
<div align="center">

|                                                                                                                       ์ด์ฐ์ : BE                                                                                                                       |                                                                                                                       ๋ณ์ฌ์ง : BE                                                                                                                       |                                                                                                                       ์ฌ์นํ : BE                                                                                                                                                                                               |
| :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: 
|                                                                                     <img width="160px" src="https://avatars.githubusercontent.com/u/119834304?v=4">                                                                                      |                                                            <img width="160px" src="https://avatars.githubusercontent.com/u/103496262?v=4" />                                                            |                                                                                     <img width="160px" src="https://avatars.githubusercontent.com/u/65848709?v=4">                                                                                      |
[@wslee220420](https://github.com/wslee220420)                                                                                                          |                                                                                                    [@wlswo](https://github.com/wlswo)                                                                                                     |                                                                                                        [@simpsonhoon](https://github.com/simpsonhoon)

</div> 

```
๐ก ๊ตฌ๊ธ ๋ก๊ทธ์ธ์ผ๋ก ์ ์  ์นํ์ ์ธ ์ง๊ฐ ์๋น์ค ๊ตฌํ
- ์ฌ์ฉ์๋ ๋๋ชจ๋ ๋จ์ด๋ฅผ ๋ณด๊ดํ์ง ์์๋ ๋ฉ๋๋ค.
- ๊ตฌ๊ธ ๋ก๊ทธ์ธ๊ณผ ๋น๋ฐ๋ฒํธ๋ฅผ ์ด์ฉํ์ฌ PrivateKey๋ฅผ ๋ณด๊ดํฉ๋๋ค.
- ์ง๊ฐ์ฃผ์์ ํธ๋์ญ์ ๋ด์ฉ์ ์ฝ๊ฒ ํธ๋ํน ํ  ์ ์์ต๋๋ค.
```

### ํต์ฌ ๊ธฐ๋ฅ

- ์ง๊ฐ
    - ๊ตฌ๊ธ ๋ก๊ทธ์ธ์ ํตํ ์ง๊ฐ ์์ฑ
    - ์ ์ ์ ๋ฉ์ผ์ฃผ์์ ๋น๋ฐ๋ฒํธ๋ก Kerstore ๊ด๋ฆฌ
    - ๋ฉํฐ ์ฒด์ธ์ ๋ณด์ ํ๋ ์ฝ์ธ/ํ ํฐ ์์ฐ์ ํ๋ฒ์ ํ์ธ
    - ์ฝ์ธ, ํ ํฐ ์ ์ก
    - ๋ฉํฐ์๊ทธ
- ๋ฐ๋ชฌ
    - ๋ฉํฐ ์ฒด์ธ ํธ๋์ญ์ ๋ทฐ
    - ๊ฐ ์ฒด์ธ๋ณ ํธ๋์ญ์ ์ด๋ฒคํธ ๊ฐ์ง

### ์ฐจ๋ณ์ 

- ์ ์ ๊ฐ ๋๋ชจ๋ ๋จ์ด๋, PrivateKey๊ด๋ฆฌ๋ฅผ ํ์ง ์์๋ ๋ฉ๋๋ค.
- ์๋น์ค์์ ๊ฒ์ฆํ ํ ํฐ์ ์ปจํธ๋ํธ๋ฅผ ์ถ๊ฐ ํ์ง ์์๋ ํ์ธ์ด ๊ฐ๋ฅํฉ๋๋ค.
- ์ฝ์ธ/ํ ํฐ์ ๊ฑฐ๋๋ด์ญ์ ์ฝ๊ฒ ํ์ธํ  ์ ์์ต๋๋ค.
- ๋คํธ์ํฌ๋ฅผ ์ ํํ์ง ์์๋ ์ฌ์ฉํ  ์ ์์ต๋๋ค.

## ์ฌ์ฉ๋ฒ
- ๊ตฌ๊ธ ๋ก๊ทธ์ธ, ์ง๊ฐ ์์ฑ
    - ์ง๊ฐ ์์ฑ ๋ฒํผ ํด๋ฆญ์ ์ง๊ฐ ์์ฑ
    - <img src="https://user-images.githubusercontent.com/103496262/215489595-bb2aa824-b05f-463f-b17a-a3a9d720eff5.gif">
- ์ง๊ฐ ์์ฐ ๊ฐ์ ธ์ค๊ธฐ
    - ์ง๊ฐ ์ฃผ์๋ก ์์ฐ ํธ์ถ
        - [http://localhost:8080/wallet/balance?address=0x314613c](http://localhost:8080/wallet/balance?address=0x314613c08Cb38e3d782688e86f61a563D8959574)
        - ๊ฒฐ๊ณผ ๊ฐ
            
            ```json
            {
                "coinInfos :": [
                    {
                        "Contract": "",
                        "SymbolName": "WEMIX",
                        "BalanceOf": "7.454006382490955452",
                        "Price": 0,
                        "Network": "WEMIX"
                    },
                    {
                        "Contract": "",
                        "SymbolName": "ETH",
                        "BalanceOf": "0.783477864499520535",
                        "Price": 0,
                        "Network": "ETH"
                    },
                    {
                        "Contract": "",
                        "SymbolName": "KLAY",
                        "BalanceOf": "149.932127875",
                        "Price": 0,
                        "Network": "KLAY"
                    }
                ],
                "tokenInfos :": [
                    {
                        "Contract": "0x9Fa7F4E848Df29B3F653c47cC12b4c9bBCf2b99c",
                        "SymbolName": "WAL",
                        "BalanceOf": "7.7731027138210594724e+08",
                        "Price": 0,
                        "Network": "WEMIX"
                    },
                    {
                        "Contract": "0x718B42c6E706383DB5e9dc1C1356f417E00b3977",
                        "SymbolName": "EWAL",
                        "BalanceOf": "9.999476683190259263e+08",
                        "Price": 0,
                        "Network": "ETH"
                    },
                    {
                        "Contract": "0xbccfb43e61bc1726861055f9169b817298441070",
                        "SymbolName": "KWAL",
                        "BalanceOf": "1e+09",
                        "Price": 0,
                        "Network": "KLAY"
                    }
                ]
            }
            ```
            
- ์ฝ์ธ, ํ ํฐ ์ ์ก
    - <img src="https://user-images.githubusercontent.com/103496262/215489633-7aaac118-93d0-40d9-9abf-56398895d65d.gif">
        
- ํธ๋์ญ์ ํ์ธ
    - ์๋ฏน์ค, ํด๋ ์ดํผ, ์ด๋๋ฆฌ์ ์ฒด์ธ์์ ์์ ์ ์ฝ์ธ, ํ ํฐ์ ๊ฑฐ๋๋ด์ญ ๋ฐ ํธ๋์ญ์ ๋ฆฌ์คํธ๋ฅผ ๊ฐ์ ธ์ต๋๋ค.
    - <img src="https://user-images.githubusercontent.com/103496262/215490878-9dac65ce-c0f7-4151-949a-dc4f841c6d96.gif">

- ๋ค์ค์๋ช ์ง๊ฐ ์์ฑ๊ณผ ๊ณ์ฝ ์ ์ถ
    - <img src="https://user-images.githubusercontent.com/103496262/215489677-18be0ba6-0a8b-43be-a979-f9a3333b9860.gif">
    - <img src="https://user-images.githubusercontent.com/103496262/215490858-a195fb13-6caf-4b62-a7ee-7311d05fd86a.gif">

### ๊ฐ๋ฐ ํ๊ฒฝ
- ์ธ์ด
    - Go
- Database
    - mongodb
- API Test
    - Postman
- ํ์ํด
    - Git
    - Discord
### ์ฌ์ฉ ๊ธฐ์ 

### ERD
- ๋ฐ์ดํฐ๋ฒ ์ด์ค ์ค๊ณ
    - ์ ์  Database
        - MEMBER
        <pre><code>
        type User struct {
         ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //๊ธฐ๋ณธํค
         Email    string             `json:"email" bson:"email"`                 //SNS ID
         Address  string             `json:"address" bson:"address"`             //์ง๊ฐ์ฃผ์
        }
        </pre></code>
        - WALLET
    - ๋ฐ๋ชฌ Database
        - WemixBlock
        - KlaytnBlock
        - EthereumBlock
        <pre><code>
        type Block struct {
         BlockHash    string        `bson:"blockHash"`
         BlockNumber  uint64        `bson:"blockNumber"`
         GasLimit     uint64        `bson:"gasLimit"`
         GasUsed      uint64        `bson:"gasUsed"`
         Time         uint64        `bson:"timestamp"`
         Nonce        uint64        `bson:"nonce"`
         Transactions []Transaction `bson:"transactions"`
        }

        type Transaction struct {
         TxHash      string `bson:"hash"`
         From        string `bson:"from"`
         To          string `bson:"to"` 
         Nonce       uint64 `bson:"nonce"`
         GasPrice    uint64 `bson:"gasPrice"`
         GasLimit    uint64 `bson:"gasLimit"`
         Amount      uint64 `bson:"amount"`
         BlockHash   string `bson:"blockHash"`
         BlockNumber uint64 `bson:"blockNumber"`
        }
        </pre></code>
        
### ์๋ฒ ์ํคํ์ฒ 
- ์๋ฒ๊ตฌ์ฑ
    - Application Server 1๋
    - DataBase Server 1๋
    - Deamon Server 3๋
        
<br>

### API ๋ช์ธ์
---
๋ก๊ทธ์ธ ๊ด๋ จ
```
/auth/google/login                  [GET]        // @Description  ๊ตฌ๊ธ ๋ก๊ทธ์ธ
/auth/google/callback               [GET]        // @Description  ๊ตฌ๊ธ ๋ก๊ทธ์ธ ์ฝ๋ฐฑ 
```

์ง๊ฐ ๊ด๋ จ
```
/wallet/trackAddress/{from}         [GET]        // @Description  ํน์  ์ฃผ์ ๋ฐ์ํ ํธ๋์ญ์ ๊ฐ์ ธ์ค๊ธฐ
/wallet/trackContract               [POST]       // @Description  ์ฝ์ธ  ๊ฐ์ ธ์ค๊ธฐ
/wallet/mnemonics                   [POST]       // @Description  ๋๋ชจ๋ ์์ฑ
/wallet/                            [POST]       // @Description  ์ง๊ฐ ์์ฑ
/wallet/balance                     [GET]        // @Description  ์์ฐ ์ ๋ณด ๊ฐ์ ธ์ค๊ธฐ
/wallet/transfer                    [POST]       // @Description  ์ฝ์ธ,ํ ํฐ ์ ์ก
``` 

๋ค์ค์๋ช์ง๊ฐ ๊ด๋ จ
```
/multisigwallet/                    [GET]         // @Description ๋ค์ค์๋ช์ง๊ฐ ์์ฑ ํ์ด์ง ๋ฐํ
/multisigwallet/                    [POST]         // @Description ๋ค์ค์๋ช์ง๊ฐ์์ฑ
/multisigwallet/submit              [POST]         // @Description ๋ค์ค์๋ช ๊ณ์ฝ ์ ์ถ
/multisigwallet/confirm             [POST]         // @Description ๊ณ์ฝ ์ปจํ
/multisigwallet/txCount/{wallet}    [GET]         // @Description ๊ณ์ฝ ๊ฐ์ ๊ฐ์ ธ์ค๊ธฐ
/multisigwallet/owners/{wallet}     [GET]         // @Description ์ง๊ฐ ์์ ์ ๋ชฉ๋ก ๋ฐํ
```


