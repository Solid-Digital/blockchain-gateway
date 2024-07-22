const Web3 = require('web3');

const username = 'YOUR_USERNAME_HERE';
const password = 'YOUR_PASSWORD_HERE';
const url = 'YOUR_URL_HERE'; // make sure to strip the wss:// prefix!

const wsSecURL = 'wss://' + username + ':' + password + '@' + url;
const web3ws = new Web3(new Web3.providers.WebsocketProvider(wsSecURL));

const tether_abi = require('./tether_abi.json'); // you can copy the ABI for the Tether contract here: https://etherscan.io/address/0xdac17f958d2ee523a2206206994597c13d831ec7#code
const contractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"; // tether contract address, check here: https://etherscan.io/address/0xdac17f958d2ee523a2206206994597c13d831ec7

var tetherContract = new web3ws.eth.Contract(tether_abi, contractAddress);

console.log('starting listener');
tetherContract.events.allEvents()
  .on('data', (event) => {
    console.log(event)
    // store event to local database or present to client
  })
  .on('changed', (event) => {
    console.log(event)
    // i.e. remove event from local database or remove for client state
  })
  .on('error', (error) => {
    // something went wrong
    console.log(error);
  });


