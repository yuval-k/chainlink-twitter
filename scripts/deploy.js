Web3 = require('web3')
web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:8545"));
code = fs.readFileSync('Voting.sol').toString()
solc = require('solc')
compiledCode = solc.compile(code)
abiDefinition = JSON.parse(compiledCode.contracts[':LinkToken'].interface)
Contract = web3.eth.contract(abiDefinition)
byteCode = compiledCode.contracts[':LinkToken'].bytecode
deployedContract = Contract.new()
deployedContract.address