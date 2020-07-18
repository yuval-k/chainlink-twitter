<script>
  import "./main.css";
  import Web3Connector from "./Web3Connector";
  import contractData from "../../contracts/twitterconsumer/build/contracts/TwitterConsumer.json";
  import linkContractData from "../../LinkToken/build/contracts/LinkToken.json";

  let name = "";
  let web3;
  let account;
  let linkcontract;

  let link = "";
  let deadline = "86400";
  let beneficiary = "0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0";
  let amount = "1";
  let linkAmount = "1";
  let handle = "KohaviYuval";
  let text = "yes";
  let deployedContractAddress = "";

  let isReady = false;

  // changeme
  let oracle = "";
  let jobid = "";

  async function deploy() {
    try {
      let accounts = await web3.eth.getAccounts();
      let requester = accounts[0];
      console.log(
        "from:",
        requester,
        "params:",
        link,
        deadline,
        beneficiary,
        amount,
        handle,
        text,
        oracle,
        jobid
      );
      var bytecode = contractData.bytecode;

      var c = new web3.eth.Contract(contractData.abi);
      var deployment = c.deploy({
        data: bytecode,
        arguments: [
          link,
          deadline,
          beneficiary,
          amount*10e18,
          handle,
          text,
          oracle,
          jobid
        ]
      });
      console.log("estimating gas");
      let gasAmount = await deployment.estimateGas({ from: requester });
      console.log("gas estimate:");
      console.log(gasAmount);
      await deployment
        .send({ from: requester, gas: 3000000 })
        .on("receipt", function(receipt) {
          console.log("deployed contract address: ", receipt.contractAddress);
          deployedContractAddress = receipt.contractAddress;
        });
    } catch (err) {
      name = "error " + err;
    }
  }

  async function fund() {
    var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
    let txhash = await c.methods
      .fund()
      .send({ from: account, gas: 900000, value: amount*10e18});
    console.log("txhash: ", txhash);
    txhash = await linkcontract.methods
      .transfer(deployedContractAddress, linkAmount*10e18)
      .send({ from: account, gas: 900000 });
    console.log("txhash: ", txhash);
  }
  async function ready() {
    var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
    isReady = await c.methods.ready().call({ from: beneficiary });
    console.log("isReady: ", isReady);
  }
  async function requestApproval() {
    var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
    let txhash = await c.methods
      .requestApproval()
      .send({ from: account, gas: 900000 });
    console.log("txhash: ", txhash);
  }
  async function withdraw() {
    var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
    let txhash = await c.methods
      .withdraw()
      .send({ from: beneficiary, gas: 900000 });
    console.log("txhash: ", txhash);
  }
</script>

<style>
  .has-default {
    background-color: gray;
  }
  .originator {
    background-color: lightblue;
  }
  .beneficiary {
    background-color: lightpink;
  }
  .approver {
    background-color: lightgreen;
  }
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
  }

  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 4em;
    font-weight: 100;
  }
  h1 {
    @apply bg-black text-white;
  }

  .tweet {
    @apply bg-white rounded-lg p-6;
    max-width: 500px;
    margin: 5px auto 0;
  }
  .tweet .btn {
    background-color: rgb(29, 161, 242);
    font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
      Roboto, Ubuntu, "Helvetica Neue", sans-serif;
    color: rgb(255, 255, 255);
    font-size: 15px;
    font-weight: 700;
    padding-right: 15px;
    padding-left: 15px;
    border-radius: 9999px;
  }

  .step {
    padding-top: 10px;
    padding-bottom: 10px;
    margin-top: 10px;
    margin-bottom: 10px;
    border-radius: 10px;
  }

  .step-action{
    background-color: wheat;
    border-radius: 5px;
    padding-right: 5px;
    padding-left: 5px;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>

<main>

  <h1>Hello {name}!</h1>
  <Web3Connector
    bind:web3
    bind:account
    bind:link={linkcontract}
    bind:linkAddr={link} />
  <div class="step originator">
    <h2>Deploy new contract</h2>
    In this step we will deploy a new contract, and have the receipent verify that it's terms are agreeable.
    <input
      type="text"
      placeholder="link"
      bind:value={link}
      class="has-default" />
    <input
      type="text"
      placeholder="deadline"
      bind:value={deadline}
      class="has-default" />
    <input type="text" required="required" placeholder="originator" bind:value={account} />
    <input type="text" required="required" placeholder="beneficiary" bind:value={beneficiary} />
    <input type="number" min="1" required="required" placeholder="amount" bind:value={amount} />
    <input type="text" required="required" placeholder="handle" bind:value={handle} />
    <input type="text" required="required" placeholder="text" bind:value={text} />
    <input
      type="text"
      required="required"
      placeholder="oracle address"
      bind:value={oracle}
      class="has-default" />
    <input
      type="text"
      required="required"
      placeholder="job id"
      bind:value={jobid}
      class="has-default" />

    <button class="step-action" on:click={deploy}>Deploy</button>
    {#if deployedContractAddress}
      <span>deployed contract address: {deployedContractAddress}</span>
    {/if}
  </div>
  <div class="step originator">
    <h2>Fund</h2>
    Once the terms are agreed upon, the contract should be funded with ETH(for the beneficiary) and LINK
    (for the oracle).
    <input
      type="text"
      placeholder="deployedContractAddress"
      bind:value={deployedContractAddress}
      class="has-default" />
    <input
      type="text"
      placeholder="account"
      bind:value={account}
      class="has-default" />
    <input type="number" min="1" placeholder="amount" bind:value={amount} />
    <input type="number" min="1" placeholder="linkAmount" bind:value={linkAmount} />

    <button class="step-action" on:click={fund}>Fund</button>
  </div>

  <div class="step beneficiary">
    <h2>Check Ready</h2>
    To verify that the contract is funded, the beneficiary can use this step to check that the contract
    has the ETH amount agreed upon, and the LINK amount for the oracle. Once ready, the beneficiary 
    can executre the real world transaction.
    <input
      type="text"
      placeholder="deployedContractAddress"
      bind:value={deployedContractAddress}
      class="has-default" />
    <input
      type="text"
      placeholder="beneficiary"
      bind:value={beneficiary}
      class="has-default" />

    <button class="step-action" on:click={ready}>Check Ready</button>
    <div>
      <label for="checkboxIsReady">Ready:</label>
      <input
        type="checkbox"
        id="checkboxIsReady"
        disabled="disabled"
        bind:checked={isReady} />
    </div>
  </div>

  <div class="step approver">
    <h2>Approve</h2>
    Once the real world transaction was performed, the trusted 3rd party needs to approve by publicly
    tweeting the agreed upton text. Please tweet as: @{handle} a tweet containing this: {text} Here's a quick
    link for you:
    <div class="tweet">
      <div class="text-lg md:text-left">@{handle}</div>
      <div>{text}</div>
      <div class="md:text-right">
        <a
          class="btn"
          target="_blank"
          rel="noopener noreferrer"
          href="https://twitter.com/intent/tweet?ref_src=twsrc%5Etfw&amp;text={text}&amp;tw_p=tweetbutton&amp">
          <i />
          <span class="label" id="l">Tweet</span>
        </a>
      </div>
    </div>

  </div>

  <div class="step beneficiary">
    <h2>check approval</h2>
    Check that the contract is approved - the oracle did its job!
    <input
      type="text"
      placeholder="deployedContractAddress"
      bind:value={deployedContractAddress}
      class="has-default" />
    <button class="step-action" on:click={requestApproval}>Request Approval</button>
  </div>

  <div class="step beneficiary">
    <h2>withdraw</h2>
    The transaction is complete - the beneficiary can now width the funds!
    <input type="text" placeholder="beneficiary" bind:value={beneficiary} />
    <input
      type="text"
      placeholder="deployedContractAddress"
      bind:value={deployedContractAddress}
      class="has-default" />
    <button class="step-action" on:click={withdraw}>Withdraw</button>
  </div>

</main>
