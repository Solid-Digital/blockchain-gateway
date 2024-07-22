pragma solidity ^0.5.1;

contract SimpleStorage {
    uint storedData;
    uint foo;
    string bar;

    constructor(uint _foo, string memory _bar) public {
        foo = _foo;
        bar = _bar;
    }

    function set(uint x) public {
        storedData = x;
    }

    function get() public view returns (uint) {
        return storedData;
    }
}