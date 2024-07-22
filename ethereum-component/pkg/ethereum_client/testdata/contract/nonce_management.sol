pragma solidity ^0.5.8;

contract NonceManagement {
    mapping(string => string) private records;

    function set(string memory _key, string memory _value) public {
        records[_key] = _value;
    }

    function get(string memory _key) public view returns (string memory) {
        return records[_key];
    }
}
