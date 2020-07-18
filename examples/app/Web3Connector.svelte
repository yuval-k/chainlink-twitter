<script>
  import "./main.css";
  import linkContractData from "../../LinkToken/build/contracts/LinkToken.json";
  import networkData from "./info.json";
  import Web3 from "web3";
  export let web3;

  export var link = null;
  export var linkAddr = "";
  export var account = null;
  var linkBalance = 0;
  var ethBalance = 0;

  let address = "ws://127.0.0.1:32000";
  let networkInfo = null;
  let currentNetwork = "<not connected>";
  let currentNetworkName = "<not connected>";
  let currentChain = null;
  let metamaskStatus = "MetaMask not detected. please install it first";
  let metamaskInstalled = false;

  //let web3 = new Web3('ws://localhost:8546');
  if (typeof window.ethereum !== "undefined") {
    console.log("MetaMask is installed!");
    metamaskStatus = "MetaMask is installed!";
    metamaskInstalled = true;

    ethereum.autoRefreshOnNetworkChange = false;
    //ethereum.request({ method: 'eth_requestAccounts' });
    // reload on chain change to keep it simple.
    ethereum.on("chainChanged", chainId => {
      if (currentChain != null && currentChain != chainId) {
        document.location.reload();
      }
      currentChain = chainId;
    });

    ethereum.on("accountsChanged", newAccountCb);
  }

  async function newAccountCb(accounts) {
    setupInjectedWeb3();
    //  account = accounts[0];
    // do something with new account here
    // newAccount();
    connected();
  }

  async function newAccount() {
    const accounts = await web3.eth.getAccounts();
    account = accounts[0];
    // do something with new account here
  }

  export async function showbalances() {
    ethBalance = await web3.eth.getBalance(account);
    let newlinkaddr = "";
    let newnetname = "";
    if (networkInfo != null){
      newlinkaddr = networkInfo.LINK;
    }
    if (newlinkaddr == "") {
      link = null;
      linkBalance = 0;
    } else if (linkAddr != newlinkaddr) {
      linkAddr = newlinkaddr;
      link = new web3.eth.Contract(linkContractData.abi, linkAddr);
    }  
    linkBalance = await link.methods.balanceOf(account).call({ from: account });
    console.log("linkBalance", linkBalance);
  }

  function setnetworkdata(){
    for (let potentialNetwork of networkData) {
      if (potentialNetwork.chainId == currentChain) {
        networkInfo = potentialNetwork;
        return;
     }
    }
    networkInfo = null;
  }
  async function connected() {
    try {
      currentNetwork = await web3.eth.net.getId();
      currentChain = await web3.eth.getChainId();
      setnetworkdata();

    if (networkInfo != null){
      currentNetworkName = networkInfo.name;
    }

      await newAccount();
      await showbalances();
    } catch (err) {
      currentNetwork = "not connected: "+err;
    }
  }

  async function setupInjectedWeb3() {
    if (typeof window.ethereum === "undefined") {
      throw Exception("MetaMask is not installed!");
    }

    if (typeof web3 === "undefined") {
      web3 = new Web3(window.ethereum);
    } else {
      web3.setProvider(window.ethereum);
    }
    await window.ethereum.enable();
  }
  async function connectInjected() {
    await setupInjectedWeb3();
    await connected();
  }

  async function connectAddress(addr) {
    if (typeof web3 === "undefined") {
      web3 = new Web3(addr);
    } else {
      web3.setProvider(addr);
    }
    await connected();
  }

  async function handleClick() {
    try {
      await connectInjected();
    } catch {
      currentNetwork = "error";
    }
  }

  async function handleAddress() {
    try {
      await connectAddress(address);
    } catch {
      currentNetwork = "error";
    }
  }
</script>

<style>
  button{
    @apply rounded-lg bg-gray-300 shadow-lg;
  }
  .widget{
    display: grid;
    grid-template-columns: 1fr 1fr;
  }
  .twocolums{
  grid-column: 1 / 3;
  }

</style>

<div class="widget">
  <span>{metamaskStatus}</span>
  <button on:click={handleClick} disabled={!metamaskInstalled}><u>Option 1:</u> Connect with Metamask</button>
  <input type="text" placeholder="provider" bind:value={address} />
  <button on:click={handleAddress}><u>Option 2:</u> Connect to a Provider</button>
  <span>Current Network:</span><span>{currentNetworkName} ({currentNetwork},{currentChain})</span>
  <span>Current Account:</span><span>{account}</span>
  <span>LINK balance</span><span>{linkBalance/1e18}</span>
  <span>ETH balance</span><span>{ethBalance/1e18}</span>
</div>
