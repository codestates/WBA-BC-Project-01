# WBA-BC-Project-01

# 📖 목차 
 - [소개](#소개) 
 - [사용법](#사용법)
 - [디렉토리 구조](#디렉토리-구조)
 - [개발 환경](#개발-환경)
 - [사용 기술](#사용-기술)
 - [ERD](#erd)
 - [서버 아키텍처](#서버-아키텍처) 
 - [Api 명세서](#api-명세서)
 - [피드백 개선](#피드백-개선)
 - [주요 기능](#주요-기능)
    - [데몬 서버](#데몬서버)
        - [이더리움 데몬](#이더리움)
        - [클레이튼 데몬](#클레이튼)
        - [위믹스 데몬](#위믹스)
    - [멀티 체인 지갑](#멀티체인지갑)
        - [로그인](#로그인)
        - [코인전송](#코인전송)
       
    
    
# 소개 
 - 구글 소셜 로그인을 이용한 멀티 체인 지갑
 거래소가 아닌 개인지갑을 이용 예정인 신규 유저는 지갑생성부터 어려움을 느낀다. 그래서 구글 로그인을 하면 자동으로 지갑을 생성하고, 유저가 니모닉 단어와 지갑 키를 관리하지 않아도 되는 지갑이 있으면 좋겠다고 생각했습니다. 

 저희 프로젝트는 기존 소셜과 동일하게 메일주소와 패스워드만 있으면 누구나 사용 가능한 유저 친화적인 지갑을 제공합니다.
 
<aside>
💡 **구글 로그인으로 유저 친화적인 지갑 서비스 구현**

- 사용자는 니모닉 단어를 보관하지 않아도 됩니다.
- 구글 로그인과 비밀번호를 이용하여 PrivateKey를 보관합니다.
- 지갑주소의 트랜잭션 내용을 쉽게 트래킹 할 수 있습니다.
</aside>

### 핵심 기능

- 지갑
    - 구글 로그인을 통한 지갑 생성
    - 유저의 메일주소와 비밀번호로 Kerstore 관리
    - 멀티 체인의 보유하는 코인/토큰 자산을 한번에 확인
    - 코인, 토큰 전송
    - 멀티시그
- 데몬
    - 멀티 체인 트랜잭션 뷰
    - 각 체인별 트랜잭션 이벤트 감지

### 차별점

- 유저가 니모닉 단어나, PrivateKey관리를 하지 않아도 됩니다.
- 서비스에서 검증한 토큰은 컨트랙트를 추가 하지 않아도 확인이 가능합니다.
- 코인/토큰의 거래내역을 쉽게 확인할 수 있습니다.
- 네트워크를 전환하지 않아도 사용할 수 있습니다.


## 사용법
- 구글 로그인, 지갑 생성
    - 지갑 생성 버튼 클릭시 지갑 생성
        - [http://localhost:8080/auth/google/login](http://localhost:8080/auth/google/login)
    - 패스워드 입력 후 지갑 생성
        - 패스워드는 개인키를 복호화할 때 사용
    - 지갑 생성 후, 메일주소와 지갑주소 리턴
        - [http://localhost:8080/wallet/](http://localhost:8080/wallet/)
- 지갑 자산 가져오기
    - 지갑 주소로 자산 호출
        - [http://localhost:8080/wallet/balance?address=0x314613c](http://localhost:8080/wallet/balance?address=0x314613c08Cb38e3d782688e86f61a563D8959574)
        - 결과 값
            
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
            
- 코인, 토큰 전송
    - 네트워크 정보, 주소, 수량, 메일주소, 패스워드를 입력받아 실행
        - [http://localhost:8080/wallet/transfer](http://localhost:8080/wallet/transfer)
    - 결과 값
        
        ```json
        {
            "tx :": "0xd7de8be613af5fbccea08025a5258f716dbf4ee178e7342d37e05533bf68c6e8"
        }
        ```
        
- 트랜잭션 확인
    - 위믹스, 클레이튼, 이더리움 체인에서 자신의 코인, 토큰의 거래내역 및 트랜잭션 리스트를 가져온다.

### 개발 환경
- 언어
    - Go

### 사용 기술


### 디렉토리 구조 


### ERD
- 데이터베이스 설계
    - MEMBER
    - WALLET
    - BLOCK
        - WemixBlock
        - KlaytnBlock
        - EthereumBlock
        
### 서버 아키텍처 
- 서버구성
    - Application Server 1대
    - DataBase Server 1대
    - Deamon Server 3대
        
### API 명세서

# 피드백 개선

# 주요 기능

