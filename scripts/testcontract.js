const Web3 = require('web3');
const fs = require('fs');

var jsonFile = './contracts/testconsumer/build/contracts/ATestnetConsumer.json';
var parsedContract = JSON.parse(fs.readFileSync(jsonFile));
var jsonInterface = parsedContract.abi;

const web3 = new Web3('http://localhost:32000');

[_,_, contractAddress, oracle, jobid] = process.argv;
(async () => {
    let accounts = await web3.eth.getAccounts();
    let requester = accounts[0];
    console.log("from:",requester,"params:", contractAddress, oracle, jobid)
    
    var c = new web3.eth.Contract(jsonInterface, contractAddress);
    txhash = await c.methods.requestEthereumPrice(oracle, jobid).send({from: requester, gas:900000});
    console.log(txhash);
})()