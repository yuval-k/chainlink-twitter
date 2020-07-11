function request(contract, contractAddress, requesterAddress, oracle, jobid) {
    var c = web3.eth.contract(contract).at(contractAddress);
    console.log("address: ", c.address, "requesterAddress: ", requesterAddress, "oracle: ", oracle, "jobid", jobid);
    return c.requestEthereumPrice(oracle, jobid, {from: requesterAddress});
}
