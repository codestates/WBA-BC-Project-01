package services

type MultiSigWalletService interface {
	//이메일,패스워드, 주소배열, 컨펌수 를받아 멀티시그지갑 주소를 리턴합니다.
	CreateMultiSigWallet(email string, password string, addresss []string, NumConfirmRequired uint) (string, error)
}

/*
	멀티시그지갑 생성하기 - 컨트랙트 배포 (이메일,패스워드, 주소들, 컨펌수) (트랜잭션 해시, 에러)
	트랜잭션 제출하기
	제출된 트랜잭션 확인하기
	제출된 트랜잭션 거절하기
	트랜잭션 실행하기
	소유자들 반환
	현재 총 트랜잭션 수 반환
	특정 트랜잭션 상세정보 가져오기
*/
