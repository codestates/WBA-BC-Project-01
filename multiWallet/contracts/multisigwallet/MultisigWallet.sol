// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract MultiSigWallet {
    // --------------------------- 🗞 EVENT ---------------------------
    event Deposit(address indexed sender, uint amount, uint balance);
    event SubmitTransaction(
        address indexed owner,
        uint indexed txIndex,
        address indexed to,
        uint value,
        bytes data
    );
    event ConfirmTransaction(address indexed owner, uint indexed txIndex);
    event RevokeConfirmation(address indexed owner, uint indexed txIndex);
    event ExecuteTransaction(address indexed owner, uint indexed txIndex);
    // --------------------------- 🗞 EVENT ---------------------------

    // --------------------------- 🧺 VARIABLES ---------------------------
    address[] public owners;                 //소유자들 
    mapping(address => bool) public isOwner; //소유자 : boolean
    uint public numConfirmationsRequired;    //컨펌개수

    //트랜잭션 구조체 data를 바이트형태로 담을수 있는게 특징
    struct Transaction {
        address to;
        uint value;
        bytes data;
        bool executed;
        uint numConfirmations;
    }

    //[트랜잭션 인덱스][ [소유자 주소][ true/false ] ] 이중맵 자료구조로 특정 트랜잭션에대해 소유자가 컨펌 했는지 여부를 판단하기 위한 변수
    mapping(uint => mapping(address => bool)) public isConfirmed;
    
    //트랜잭션들 
    Transaction[] public transactions;
    // --------------------------- 🧺 VARIABLES ---------------------------

    //소유자만 호출 가능 
    modifier onlyOwner() {
        require(isOwner[msg.sender], "not owner");
        _;
    }

    //트랜잭션 존재 유무 검사 , revert 상황은 [0,1,2] 상황에서 5번째 트랜잭션을 조회하는 경우 
    modifier txExists(uint _txIndex) {
        require(_txIndex < transactions.length, "tx does not exist");
        _;
    }

    // 트랜잭션이 실행된 상태인지 
    modifier notExecuted(uint _txIndex) {
        require(!transactions[_txIndex].executed, "tx already executed");
        _;
    }
    // 트랜잭션이 컨펌된 상태인지 
    modifier notConfirmed(uint _txIndex) {
        require(!isConfirmed[_txIndex][msg.sender], "tx already confirmed");
        _;
    }

    constructor(address[] memory _owners, uint _numConfirmationsRequired) {
        require(_owners.length > 0, "owners required");
        require(
            _numConfirmationsRequired > 0 &&
                _numConfirmationsRequired <= _owners.length,
            "invalid number of required confirmations"
        );

        for (uint i = 0; i < _owners.length; i++) {
            address owner = _owners[i];

            require(owner != address(0), "invalid owner");
            require(!isOwner[owner], "owner not unique");

            isOwner[owner] = true;
            owners.push(owner);
        }

        numConfirmationsRequired = _numConfirmationsRequired;
    }
    //해당 지갑으로 예치할경우 즉 이 컨트랙트로 이더를 보내는 경우
    receive() external payable {
        emit Deposit(msg.sender, msg.value, address(this).balance);
    }

    //트랜잭션 제출 
    function submitTransaction(address _to,uint _value,bytes memory _data) public onlyOwner {
        //트랜잭션 인덱스 저장 
        uint txIndex = transactions.length;
        //트랜잭션 배열에 새로운 트랜잭션 저장
        transactions.push(
            Transaction({
                to: _to,
                value: _value,
                data: _data,
                executed: false,
                numConfirmations: 0
            })
        );

        emit SubmitTransaction(msg.sender, txIndex, _to, _value, _data);
    }
    //트랜잭션 컨펌 요청하기 
    //소유자만 호출가능, 존재하는 트랜잭션인지, 실행되지 않았는지, 컨펌이 안된 트랜잭션인지 검사
    function confirmTransaction(uint _txIndex ) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) notConfirmed(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        transaction.numConfirmations += 1;
        isConfirmed[_txIndex][msg.sender] = true;

        emit ConfirmTransaction(msg.sender, _txIndex);
    }

    //트랜잭션 실행 
    //트랜잭션이 존재, 실행되지 않은 상태인지 검사
    //설정한 요구 컨펌수에 충족하는지 검사 
    function executeTransaction(uint _txIndex) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];

        require(
            transaction.numConfirmations >= numConfirmationsRequired,
            "cannot execute tx"
        );

        transaction.executed = true;

        (bool success, ) = transaction.to.call{value: transaction.value}(
            transaction.data
        );
        require(success, "tx failed");

        emit ExecuteTransaction(msg.sender, _txIndex);
    }
    //트랜잭션 취소
    //소유자만 호출가능, 존재하는 트랜잭션에 대한 요청인지, 이미 실행된 트랜잭션은 아닌지 검사
    function revokeConfirmation( uint _txIndex) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        //이미 취소한 상태인지 검사
        require(isConfirmed[_txIndex][msg.sender], "tx not confirmed");
        //컨펌수 -1
        transaction.numConfirmations -= 1;
        //해당 트랜잭션의 컨펌을 취소상태로 변경
        isConfirmed[_txIndex][msg.sender] = false;
        //이벤트 기록
        emit RevokeConfirmation(msg.sender, _txIndex);
    }
    //소유자들 반환
    function getOwners() public view returns (address[] memory) {
        return owners;
    }
    //현재 까지 트랜잭션 수 반환
    function getTransactionCount() public view returns (uint) {
        return transactions.length;
    }
    //트랜잭션 상세정보 반환 
    function getTransaction( uint _txIndex )public view
        returns (
            address to,
            uint value,
            bytes memory data,
            bool executed,
            uint numConfirmations
        )
    {
        Transaction storage transaction = transactions[_txIndex];

        return (
            transaction.to,
            transaction.value,
            transaction.data,
            transaction.executed,
            transaction.numConfirmations
        );
    }
}
