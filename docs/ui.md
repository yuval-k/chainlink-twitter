To start the UI, run:
```bash
npm run demo-dev
```

OR open the pre-built version at [examples/app/public/index.html](examples/app/public/index.html).

## Connect to the chain
Assuming you are running the setup script, click on "Option 2: Connect to a Provider". If you setup
MetaMask properly (by adding our test network to MetaMask), you can also use MetaMask (Option 1).

## Deploy contract
The first step is to deploy the contract. The fields should be pre-populated with demo defaults.
Adjust the other fields to your liking, and deploy the contract using the "Deploy" button.

## Go through the flow...

Once deployed, you can go over the flow step by step. Note that Blue steps are steps taken by the originator. Pink steps are by the beneficiary. And the green step is taken by the trusted 3rd party.

In the real world each individual will have it's own screen with just the steps relevant to him. For simplicity and to make it easier to understand the steps in order, I've combined all the steps to a single page.

Note: If MetaMask fails deploy the contract, go to Settings > Advanced > Reset Account. See this
for more details see here : https://github.com/MetaMask/metamask-extension/issues/1999