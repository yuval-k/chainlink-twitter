<script>
  import "./main.css";
  import Web3Connector from "./Web3Connector";
  import contractData from "../../contracts/twitterconsumer/build/contracts/TwitterConsumer.json";
  import linkContractData from "../../LinkToken/build/contracts/LinkToken.json";

  let name = "";
  let web3;
  let showbalances;
  let account;
  let linkcontract;

  let link = "";
  let deadline = "86400";
  let amount = "1";
  let linkAmount = "1";
  let handle = "KohaviYuval";
  let text = "yes";
  let deployedContractAddress = "";

  let fundResult = "";
  let requestApprovalResult = "";
  let withdrawResult = "";

  const eighteendecimals = "000000000000000000";

  let isReady = null;
  let isApproved = null;

  // changeme
  let beneficiary = "0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0";
  let oracle = "";
  let jobid = "";

  async function deploy() {
    try {
      let accounts = await web3.eth.getAccounts();
      let requester = accounts[0];
      let strAmount = String(amount)+eighteendecimals;
      console.log(
        "from:",
        requester,
        "params:",
        link,
        deadline,
        beneficiary,
        strAmount,
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
          strAmount,
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
    try {
      let strAmount = String(amount)+eighteendecimals;
      let strlLinkAmount = String(linkAmount)+eighteendecimals;
      var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
      let txhash = await c.methods
        .fund()
        .send({ from: account, gas: 900000, value: strAmount});
      console.log("txhash: ", txhash);
      txhash = await linkcontract.methods
        .transfer(deployedContractAddress, strlLinkAmount)
        .send({ from: account, gas: 900000 });
      console.log("txhash: ", txhash);
      await showbalances();
      fundResult = "success!";
    }catch (e){
      fundResult = "error"+ e.message;
    }
  }
  async function ready() {
    var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
    isReady = await c.methods.ready().call({ from: beneficiary });
    console.log("isReady: ", isReady);
  }
  async function requestApproval() {
    try {
      var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
      let txhash = await c.methods
        .requestApproval()
        .send({ from: account, gas: 900000 });
      console.log("txhash: ", txhash);
      requestApprovalResult = "approval requested!";
    }catch (e){
      requestApprovalResult = "error"+ e.message;
    }
  }
  async function checkIsApproved() {
    var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
    isApproved = await c.methods.approved().call({ from: beneficiary });
    console.log("isApproved: ", isApproved);
  }
  async function withdraw() {
    try {
      var c = new web3.eth.Contract(contractData.abi, deployedContractAddress);
      let txhash = await c.methods
        .withdraw()
        .send({ from: beneficiary, gas: 900000 });
      console.log("txhash: ", txhash);
      showbalances();
      withdrawResult = "withdraw success!";
    }catch (e){
      withdrawResult = "error"+ e.message;
    }
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
  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 4em;
    font-weight: 100;
  }
  .title{
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
    text-align: left;
  }
  .step h2 {
    font-weight: bold;
    text-decoration: underline;
    text-align: center;
    margin-bottom: 10px;
  }

  .step-action{
    background-color: wheat;
    border-radius: 5px;
    padding-right: 5px;
    padding-left: 5px;
    display: block;
    margin: 5px auto 0;
  }

  main {
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>

<main>

  <h1 class="title">Hello {name}!</h1>
  <section>
<h1>Connect to the blockchain</h1>
  <Web3Connector
    bind:web3
    bind:showbalances
    bind:account
    bind:link={linkcontract}
    bind:linkAddr={link} />


  </section>
  <div>
<h1>Interact with the blockchain</h1>

Note:
<ul class="list-disc">
  <li>Blue steps are performed by the originator</li>
  <li>Red steps are performed by the beneficiary</li>
  <li>Green steps are performed by the trusted 3rd party</li>
</ul>
  </div>
<section>
<div>First, create a contract</div>

  <div class="step originator">
    <h2>Deploy new contract</h2>
    In this step we will deploy a new contract, and have the receipent verify that it's terms are agreeable.
    Please customize the inputs below:
    <label for="linkErc20">the address of the LINK ERC20 coin</label>
    <input
     id="linkErc20"
      type="text"
      placeholder="link"
      bind:value={link}
      class="has-default" />
    <label for="exirtyInSeconds">the number of seconds until the contract expires in</label>
    <input
      type="number"
      min="1"
      id="exirtyInSeconds"
      placeholder="expiry"
      bind:value={deadline}
      class="has-default" />
      (that's in {Math.floor(deadline/(24*3600))} days,
       {Math.floor((deadline%(24*3600))/3600) } hours, 
        {Math.floor((deadline%3600)/60) } minutes, 
        {deadline%60 } seconds.
       )
    <label for="originatorAddress">the address of the originator; i.e. the person paying</label>
    <input type="text" id="originatorAddress" required="required" placeholder="originator" bind:value={account} />
    <label for="beneficiaryAddress">the address of the beneficiary; i.e. the person receiving the funds</label>
    <input type="text" id="beneficiaryAddress" required="required" placeholder="beneficiary" bind:value={beneficiary} />
    <label for="ethAmount">the amount of ETH the originator will pay</label>
    <input type="number" id="ethAmount" min="1" required="required" placeholder="amount" bind:value={amount} />
    <label for="twitterHandle">the twitter handle of a trusted third party, that will verify the transaction</label>
    <input type="text" id="twitterHandle" required="required" placeholder="handle" bind:value={handle} />
    <label for="approvalText">the text that needs to be tweeted, so the transaction would be approved</label>
    <input type="text" id="approvalText" required="required" placeholder="text" bind:value={text} />
    <label for="oracleAddress">oracle address</label>
    <input
      type="text"
      id="oracleAddress"
      required="required"
      placeholder="oracle address"
      bind:value={oracle}
      class="has-default" />
    <label for="oracleAddress">job ID</label>
    <input
      type="text"
      id="jobId"
      required="required"
      placeholder="job id"
      bind:value={jobid}
      class="has-default" />

    <button class="step-action" on:click={deploy}>Deploy</button>

  </div>
  </section>
  <div>
    {#if deployedContractAddress}
      <span>Deployed contract address: {deployedContractAddress}</span>
    {:else}
      Once the contract is deployed, you can continue.
    {/if}
  </div>
<section >
  <div class="step originator">
    <h2>Fund</h2>
    Once the terms are agreed upon, the contract should be funded with ETH (for the beneficiary) and LINK
    (for the oracle).

    <span>Add</span> {amount} ETH, and 
    <input type="number" min="1" placeholder="linkAmount" bind:value={linkAmount} />
    LINK from {account}.

    <button class="step-action" on:click={fund}>Fund</button>

      <div>{fundResult}</div>
  </div>

  <div class="step beneficiary">
    <h2>Check Ready</h2>
    To verify that the contract is funded, the beneficiary ({beneficiary}) can use this step to check that the contract
    has the ETH amount agreed upon, and the LINK amount for the oracle. Once ready, the beneficiary 
    can executre the real world transaction.

    <button class="step-action" on:click={ready}>Check Ready</button>
    <div>
      <label for="checkboxIsReady">Ready:</label>

    {#if isReady !== null}
      {#if isReady}
        <span>Contract is ready!</span>
      {:else}
        <span>Contract is not ready!</span>
      {/if}
    {/if}
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
    <h2>Request Oracle Approval</h2>
    The beneficiary ({beneficiary}) can now ask to check if the contract is approved - this step 
    will invoke the oracle, that will check twitter to verify the contract.

    <button class="step-action" on:click={requestApproval}>Request Oracle Approval</button>
    <div>{requestApprovalResult}</div>
  </div>

  <div class="step beneficiary">
    <h2>Check Approved</h2>
    Optional: Verify that the oracle did its job approved the contract.

    <button class="step-action" on:click={checkIsApproved}>Check Approved</button>
    <div>
      <label for="checkboxApproved">Approved:</label>

    {#if isApproved !== null}
      {#if isApproved}
        <span>Contract is approved!</span>
      {:else}
        <span>Contract is not approved!</span>
      {/if}
    {/if}
    </div>
  </div>

  <div class="step beneficiary">
    <h2>Withdraw</h2>
    The transaction is complete - the beneficiary ({beneficiary}) can now width the funds!
    <button class="step-action" on:click={withdraw}>Withdraw</button>
    <div>{withdrawResult}</div>
  </div>
</section>
</main>
