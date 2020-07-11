function request(contract, contractAddress, oracle, jobid) {
    var c = web3.eth.contract(contract).at(contractAddress);
    console.log(c.address, oracle, jobid);
    return c.requestEthereumPrice(oracle, jobid);
}
