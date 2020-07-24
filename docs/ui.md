To start the UI, run:
```bash
npm run demo-dev
```

OR open the pre-built version at [examples/app/public/index.html](examples/app/public/index.html).

## Connect to the chain
Assuming you are running the setup script, click on "Option 2: Connect to a Provider"

## Deploy contract
Add the `ORACLE_ADDR` and the `TWITTER_JOB_ID` in their perspective fields under 
the firs step ("Deploy new contract").
Adjust the other fields to your liking, and deploy the contract using the "Deploy" button.

## Go through the flow...

Once deployed, you can go over the flow step by step. Note that Blue steps are steps taken by the originator. Pink steps are by the beneficiary. And the green step is taken by the trusted 3rd party.

In the real world each individual will have it's own screen with just the steps relevant to him. For simplicity and to make it easier to understand the steps in order, I've combined all the steps to a single page.

If metamask fails deploy the contract, go to Settings > Advanced > Reset Account