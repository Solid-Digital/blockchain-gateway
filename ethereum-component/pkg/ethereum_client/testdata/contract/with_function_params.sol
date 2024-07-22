pragma solidity ^0.5.1;

contract SimpleStorage {
    function get(string memory _myString) public view returns (string memory) {
        return _myString;
    }
}