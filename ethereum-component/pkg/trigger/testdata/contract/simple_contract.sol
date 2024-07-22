pragma solidity ^0.5.1;

contract SimpleStorage {
    uint storedData;

    event Set(address indexed sender, uint256 x);

    function set() public {
        emit Set(msg.sender, 10);
    }
}