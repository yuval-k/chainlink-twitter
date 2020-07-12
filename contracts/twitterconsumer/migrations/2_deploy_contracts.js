const TwitterConsumer = artifacts.require("TwitterConsumer");

module.exports = async function(deployer) {
  let link = process.env.LINK_TOKEN;
  if (!link){
    link = "0x5b1869d9a4c187f2eaa108f3062412ecf0526b24";
  }
  await deployer.deploy(TwitterConsumer, link);
  let instance = await TwitterConsumer.deployed();
  console.log(`contract-address\t${instance.address}`)
};
