const Web3 = require('web3');
const fs = require('fs');

var jsonFile = './contracts/twitterconsumer/build/contracts/TwitterConsumer.json';
var parsedContract = JSON.parse(fs.readFileSync(jsonFile));
var jsonInterface = parsedContract.abi;
var bytecode = parsedContract.bytecode;

const web3 = new Web3('http://localhost:32000');

[_,_, link, deadline, beneficiary, amount, handle, text, oracle, 
    jobid] = process.argv;
(async () => {
    let accounts = await web3.eth.getAccounts();
    let requester = accounts[0];
    console.log("from:",requester,"params:", link, deadline, beneficiary, amount, handle, text, oracle, 
    jobid)
    
    var c = new web3.eth.Contract(jsonInterface);
    await c.deploy({data: bytecode, arguments:[link, deadline, beneficiary, amount, handle, text, oracle, 
        jobid]}).send(
        {from: requester, gas:3000000}).on('receipt', function(receipt){
            console.log("contract address: ", receipt.contractAddress);
         });
})()
