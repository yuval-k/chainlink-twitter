const Web3 = require('web3');
const fs = require('fs');

var jsonFile = './contracts/twitterconsumer/build/contracts/TwitterConsumer.json';
var parsedContract = JSON.parse(fs.readFileSync(jsonFile));
var jsonInterface = parsedContract.abi;
var bytecode = parsedContract.bytecode;

const web3 = new Web3('http://localhost:32000');

[_,_, contractaddress, beneficiary] = process.argv;
(async () => {
    console.log("from/to:",beneficiary,"params:", contractaddress)
    var c = new web3.eth.Contract(jsonInterface, contractaddress);
    let res = await c.methods.withdraw().send({from: beneficiary, gas:3000000});
    console.log(res);
})()
