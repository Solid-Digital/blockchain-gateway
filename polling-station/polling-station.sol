
pragma solidity ^0.4.24;
contract PollingStation {

    uint40 internal Summer;
    uint40 internal Winter;
    uint40 internal Keep;

    // address[] internal voters;
    mapping (uint64 => bool) private voters;

    event VoteCasted(
        address indexed pollingStation,
        uint64 indexed voterID
    );

    // The voterID is a hash of the 
    function vote(uint8 option, uint64 voterID) public returns (uint) {

        // the msg.sender is the address
        // address voterID = msg.sender;

        require(option < 3, "invalid vote casted");

        if (voters[voterID]){
            revert("This voterID has already voted");
        }

        if (option == 0) {
            Summer += 1;
        } else if (option == 1) {
            Winter += 1;
        } else if (option == 2) {
            Keep += 1;
        }

        voters[voterID] = true;
        // uint position = voters.push(voterID);
        // return position;

        emit VoteCasted(msg.sender, voterID);
        return 0;
    }

    function getVotes() public view returns (uint40 summer, uint40 winter, uint40 keep) {
        return (Summer, Winter, Keep);
    }
}