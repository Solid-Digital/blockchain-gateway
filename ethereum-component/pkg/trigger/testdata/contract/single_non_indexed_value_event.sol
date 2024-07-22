pragma solidity ^0.5.1;

contract SimpleStorage {
    event Set(address sender);

    function set() public {
        emit Set(msg.sender);
    }
}