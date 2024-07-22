pragma solidity >=0.4.22 <0.6.0;

contract Mortal {
    /* Define variable owner of the type address */
    address owner;

    /* Function to recover the funds on the contract */
    function kill() public { if (msg.sender == owner) selfdestruct(msg.sender); }
}

contract Greeter is Mortal {
    /* Define variable greeting of the type string */
    string greeting;

    /* Main function */
    function greet() public view returns (string memory) {
        return greeting;
    }
}