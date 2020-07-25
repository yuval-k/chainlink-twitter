# Problem
The goal of this project is to allow off chain service between two people to take place, and result with on-chain payment of ETH.

We'll call the person requesting the service the "originator" and the person providing the service and receiving the funds the "beneficiary".

The question is, how do we guarantee that an off-chain service has been performed?

# Solution
To ensure that funds will only be moved
once the transaction is complete and the service was performed, we use a trusted 3rd party that both originator and beneficiary trust.
Only once the trusted 3rd party approves the contract, the contract will allow withdrawing the funds 
to the beneficiary. 
If the contract is not approved within the given time-frame, the contract allows
the originator to get a refund.

As the approval is a off-chain real world event, we need a way to see that approval was indeed given.
That's where Twitter comes in! We consider a contract approved, if a specified twitter handle (representing the trusted 3rd party) tweets
a specific phrase (a.k.a approval text) provided to the contract at its creation.

A custom chainlink adapter checks twitter and approves the contract.
Once the contract is approved, the beneficiary can withdraw the funds.

# Use cases
One example use case is hiring a social media influencer to promote your brand. Using this solution, the a brand that wants twitter exposure (the originator),
create a contract with the influencer's (the beneficiary) ETH address, twitter handle and the text that should be twitter. 
As soon as the influencer tweets the agreed text, he will be able to withdraw the funds automatically from the contract.

Another use-case that comes to mind as a future expansion for this project, is corporate approval flows:
A person in the company contracts a vendor, but the transaction needs a CEO approval. The CEO can approve the contract via email to a special email address, that a chainlink adapter will listen on. (this would require writing an email adapter in addition to a twitter adapter)

# How it works?
The flow is as follows:
1. The originator creates the contract. Both parties can see terms of the contract: 
   - The identity (twitter handle) of the approver, and the approval text.
   - The amount.
   - The orcale / job used to verify
   - The expiration date
1. In the real world, the beneficiary confirms that the details of the contract are acceptable.
1. The originator funds the contract with ETH (calls `fund()` with the amount agreed upon).
1. The originator funds the contract with LINK (contract uses 1 link to pay the oracle).
1. The beneficiary can check that the contract is ready using `ready()`.
1. Now that the contract is ready, the real world transaction can happen.
1. Once completed, the trusted 3rd party tweets the approval text.
1. Someone (originator or beneficiary) calls `requestApproval()`. This triggers an oracle request.
1. Chainlink magic happens, during which the node uses twitter API to see if the approval text was tweeted.
1. Assuming the transaction was approved, the beneficiary can now call `withdraw()` to receive the funds.
1. If the contract expires and no approval was given, the originator can request a `refund()`.


# Directory structure:
```
adapter - the adapter that confirms a contract in response to a tweet mentioning it's address.
contracts - solidity contracts for oracle and consumers are here.
docker - docker files for 3rd party program that are not available in docker hub. i.e. ganache-cli.
LinkToken - git sub-module to https://github.com/smartcontractkit/LinkToken. Used to deploy the Link token to our local test chain.
examples - The UI code to for the UI demo.
manifests - kubernetes manifests to install everything automatically.
public - Skeleton files for the UI demo.
scripts - Helper web3 JS scripts used in this readme/demos to perform operations on the blockchain.
```

# Demo

We have two demos in this repo, a UI one, and a command line one.
Both require some setup. see the [setup](./docs/setup_local_testnet.md) for more info. there's a convenience [setup.sh](./setup.sh) script that you can just run (assuming you have all the command line tools in place).

Once demo-env setup is done, you can use either the [cli demo instructions](./docs/cli.md) or the [ui demo instructions](./docs/ui.md).

# FAQ

## Why was kubernetes used? 
Kubernetes is very power in the sense that the same manifests that were used
locally, can be used in production (with minor tweaks). Using kubernetes and docker allows fully
automating all the installation steps to simple `make` targets that work reliably.

## Why not just have the approver call-out to the contract to approve it?
The thought was that we want to keep blockchain interactions to a minimum, to allow
existing real-world workflows to be used on the chain gradually.