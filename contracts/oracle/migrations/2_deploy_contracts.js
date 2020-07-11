const Oracle = artifacts.require("Oracle");

module.exports = async function(deployer) {
  let link = process.env.LINK_TOKEN;
  if (!link){
    link = "0x5b1869d9a4c187f2eaa108f3062412ecf0526b24";
  }
  let nodeaddr = process.env.NODE_ADDR;
  await deployer.deploy(Oracle, link);
  let instance = await Oracle.deployed();
  await instance.setFulfillmentPermission(nodeaddr, true);
  console.log(`contract-address\t${instance.address}`)
};
