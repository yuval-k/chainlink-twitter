const Web3 = require('web3');
const fs = require('fs');

var jsonFile = './contracts/twitterconsumer/build/contracts/TwitterConsumer.json';
var parsedContract = JSON.parse(fs.readFileSync(jsonFile));
var jsonInterface = parsedContract.abi;
var bytecode = parsedContract.bytecode;

const web3 = new Web3('http://localhost:32000');

[_,_, contractaddress, amount] = process.argv;
(async () => {
    let accounts = await web3.eth.getAccounts();
    let requester = accounts[0];
    console.log("from:",requester,"params:", contractaddress, amount)
    
    var c = new web3.eth.Contract(jsonInterface, contractaddress);
    txhash = await c.methods.fund().send({from: requester, gas:900000, value: 1000000000000000000});
    console.log(txhash);
})()
