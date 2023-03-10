// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract MultiSigWallet {
    // --------------------------- π EVENT ---------------------------
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
    // --------------------------- π EVENT ---------------------------

    // --------------------------- π§Ί VARIABLES ---------------------------
    address[] public owners;                 //μμ μλ€ 
    mapping(address => bool) public isOwner; //μμ μ : boolean
    uint public numConfirmationsRequired;    //μ»¨νκ°μ

    //νΈλμ­μ κ΅¬μ‘°μ²΄ dataλ₯Ό λ°μ΄νΈννλ‘ λ΄μμ μλκ² νΉμ§
    struct Transaction {
        address to;
        uint value;
        bytes data;
        bool executed;
        uint numConfirmations;
    }

    //[νΈλμ­μ μΈλ±μ€][ [μμ μ μ£Όμ][ true/false ] ] μ΄μ€λ§΅ μλ£κ΅¬μ‘°λ‘ νΉμ  νΈλμ­μμλν΄ μμ μκ° μ»¨ν νλμ§ μ¬λΆλ₯Ό νλ¨νκΈ° μν λ³μ
    mapping(uint => mapping(address => bool)) public isConfirmed;
    
    //νΈλμ­μλ€ 
    Transaction[] public transactions;
    // --------------------------- π§Ί VARIABLES ---------------------------

    //μμ μλ§ νΈμΆ κ°λ₯ 
    modifier onlyOwner() {
        require(isOwner[msg.sender], "not owner");
        _;
    }

    //νΈλμ­μ μ‘΄μ¬ μ λ¬΄ κ²μ¬ , revert μν©μ [0,1,2] μν©μμ 5λ²μ§Έ νΈλμ­μμ μ‘°ννλ κ²½μ° 
    modifier txExists(uint _txIndex) {
        require(_txIndex < transactions.length, "tx does not exist");
        _;
    }

    // νΈλμ­μμ΄ μ€νλ μνμΈμ§ 
    modifier notExecuted(uint _txIndex) {
        require(!transactions[_txIndex].executed, "tx already executed");
        _;
    }
    // νΈλμ­μμ΄ μ»¨νλ μνμΈμ§ 
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
    //ν΄λΉ μ§κ°μΌλ‘ μμΉν κ²½μ° μ¦ μ΄ μ»¨νΈλνΈλ‘ μ΄λλ₯Ό λ³΄λ΄λ κ²½μ°
    receive() external payable {
        emit Deposit(msg.sender, msg.value, address(this).balance);
    }

    //νΈλμ­μ μ μΆ 
    function submitTransaction(address _to,uint _value,bytes memory _data) public onlyOwner {
        //νΈλμ­μ μΈλ±μ€ μ μ₯ 
        uint txIndex = transactions.length;
        //νΈλμ­μ λ°°μ΄μ μλ‘μ΄ νΈλμ­μ μ μ₯
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
    //νΈλμ­μ μ»¨ν μμ²­νκΈ° 
    //μμ μλ§ νΈμΆκ°λ₯, μ‘΄μ¬νλ νΈλμ­μμΈμ§, μ€νλμ§ μμλμ§, μ»¨νμ΄ μλ νΈλμ­μμΈμ§ κ²μ¬
    function confirmTransaction(uint _txIndex ) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) notConfirmed(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        transaction.numConfirmations += 1;
        isConfirmed[_txIndex][msg.sender] = true;

        emit ConfirmTransaction(msg.sender, _txIndex);
    }

    //νΈλμ­μ μ€ν 
    //νΈλμ­μμ΄ μ‘΄μ¬, μ€νλμ§ μμ μνμΈμ§ κ²μ¬
    //μ€μ ν μκ΅¬ μ»¨νμμ μΆ©μ‘±νλμ§ κ²μ¬ 
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
    //νΈλμ­μ μ·¨μ
    //μμ μλ§ νΈμΆκ°λ₯, μ‘΄μ¬νλ νΈλμ­μμ λν μμ²­μΈμ§, μ΄λ―Έ μ€νλ νΈλμ­μμ μλμ§ κ²μ¬
    function revokeConfirmation( uint _txIndex) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        //μ΄λ―Έ μ·¨μν μνμΈμ§ κ²μ¬
        require(isConfirmed[_txIndex][msg.sender], "tx not confirmed");
        //μ»¨νμ -1
        transaction.numConfirmations -= 1;
        //ν΄λΉ νΈλμ­μμ μ»¨νμ μ·¨μμνλ‘ λ³κ²½
        isConfirmed[_txIndex][msg.sender] = false;
        //μ΄λ²€νΈ κΈ°λ‘
        emit RevokeConfirmation(msg.sender, _txIndex);
    }
    //μμ μλ€ λ°ν
    function getOwners() public view returns (address[] memory) {
        return owners;
    }
    //νμ¬ κΉμ§ νΈλμ­μ μ λ°ν
    function getTransactionCount() public view returns (uint) {
        return transactions.length;
    }
    //νΈλμ­μ μμΈμ λ³΄ λ°ν 
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
