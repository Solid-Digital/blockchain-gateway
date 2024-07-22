pragma solidity ^0.5.1;

contract SimpleStorage {
    event Set(address indexed sender);

    function set() public {
        emit Set(msg.sender);
    }
}