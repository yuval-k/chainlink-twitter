const ATestnetConsumer = artifacts.require("ATestnetConsumer");

module.exports = async function(deployer) {
  let link = process.env.LINK_TOKEN;
  if (!link){
    link = "0x5b1869d9a4c187f2eaa108f3062412ecf0526b24";
  }
  await deployer.deploy(ATestnetConsumer, link);
  let instance = await ATestnetConsumer.deployed();
  console.log(`contract-address\t${instance.address}`)
};
