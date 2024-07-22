pragma solidity ^0.5.1;

contract SimpleStorage {
    uint storedData = 0;

    function increment() public {
        storedData++;
    }
}