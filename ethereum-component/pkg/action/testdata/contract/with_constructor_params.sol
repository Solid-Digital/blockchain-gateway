pragma solidity ^0.5.1;

contract SimpleStorage {
    uint storedData;

    constructor(uint _foo) public {
        storedData = _foo;
    }

    function set(uint x) public {
        storedData = x;
    }

    function get(uint _myInt) public view returns (uint) {
        return _myInt;
    }
}