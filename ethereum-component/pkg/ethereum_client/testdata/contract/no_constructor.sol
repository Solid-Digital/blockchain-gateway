pragma solidity ^0.5.1;

contract SimpleStorage {
    uint storedData = 123;

    function set(uint x) public {
        storedData = x;
    }

    function get() public view returns (uint) {
        return storedData;
    }
}