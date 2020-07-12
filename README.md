npm install --save ganache-cli web3 solc

go to chainlink source and do `make install` to build the contracts


To generate the db yaml, use the following:
```
helm template release bitnami/postgresql > manifests/postgresql.yaml
```

Deploy the infastrcture, starting with the ganache and the coin:
we use ganache in deterministic mode, so ADDRESS and KEY should be the same everytime. you can see them in the ganache log output.

```bash
make kind-start
make deploy-testnet
kubectl rollout status deploy/ganache
# may need to sleep here to see logs
kubectl logs deploy/ganache
make deploy-token
```

you will see:
  LinkToken: 0x5b1869d9a4c187f2eaa108f3062412ecf0526b24


Deploy chainlink node:

```bash
export ADDRESS=0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1
export KEY=0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
export LINK_TOKEN=0x5b1869d9a4c187f2eaa108f3062412ecf0526b24
# verify that Link token 0x5b1869d9a4c187f2eaa108f3062412ecf0526b24
make deploy-node

#wait for node:
kubectl rollout status deploy/chainlink

# here too, sleep might help. if this comes empty, try again after a few seconds
export NODE_ADDR=$(kubectl logs deploy/chainlink|grep "please deposit ETH into your address:"| tr ' ' '\n'|grep 0x)
```

Fund the node with ETH and LINK
```bash
# add eth to the node
geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "'${ADDRESS}'",to: "'${NODE_ADDR}'", value: "10000000000000000000"})'

# add link to the node
geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("fund.js");transfer("'$LINK_TOKEN'", "'$ADDRESS'", "'$NODE_ADDR'");'

# to verify (optional), check node balance:
geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("fund.js");getbalance("'$LINK_TOKEN'", "'$NODE_ADDR'");'
```

Deploy oracle:

```bash
# this will also add the node to the oracle (by using the address in the env-var )
npm run deploy-oracle | tee node-tmp.txt
export ORACLE_ADDR=$(grep "contract-address" node-tmp.txt | cut -f 2)
rm node-tmp.txt
```

This part is silly, but I'm not sure how to setup api keys otherwise.
port forward to the ui:
```bash
kubectl port-forward deploy/chainlink 6688&
```
Log-in in the UI (foo@example.com/apipassword) and then open a js console and paste:

```js
response = await fetch("/v2/user/token", {
  method: "post",
  headers: {
    "Accept": "application/json",
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    "password": "apipassword"
  })
});
body = await response.json();
console.log("export ACCESS_KEY="+body.data.attributes.accessKey);
console.log("export SECRET_KEY="+body.data.attributes.secret);
```

you should see output similar to this:
```bash
export ACCESS_KEY=65505b0d8c6d4ef7a9566889087e1634
export SECRET_KEY=/cFPLUGDi0iDi4aXnS+CsxLWNXhRTOkQXxO5Ne8Au+kLqKhjRz4QMXHL9nejlPtb
```

Paste it in the terminal before continuing. Now we can use the API keys to create the
jobs:

```bash
# optional: verify that the node sees its balances:
# curl http://localhost:6688/v2/user/balances -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY"

# create the jobs:
curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"ethbytes32"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httppost"},{"type":"jsonparse"},{"type":"ethbytes32"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"multiply"},{"type":"ethint256"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"multiply"},{"type":"ethuint256"},{"type":"ethtx"}]}'


# save job id of EthUint256
JOB_ID=$(curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"ethbool"},{"type":"ethtx"}]}' | jq .data.id -r)
```

Using the node!

create a consumer:

```bash
npm run deploy-testconsumer | tee node-tmp.txt
export TEST_CONTRACT_ADDR=$(grep "contract-address" node-tmp.txt | cut -f 2)
rm node-tmp.txt
# fund the contact:
geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("fund.js");transfer("'$LINK_TOKEN'", "'$ADDRESS'", "'$TEST_CONTRACT_ADDR'");'
# verify funds contact:
geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("fund.js");getbalance("'$LINK_TOKEN'", "'$TEST_CONTRACT_ADDR'");'

# create a script with our test contract
echo "var contract_ = " $(cat contracts/testconsumer/build/contracts/ATestnetConsumer.json|jq .abi) ";function contract() {return contract_;} " > scripts/contract.js

geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("contract.js");loadScript("request.js");request(contract(), "'$TEST_CONTRACT_ADDR'", "'$ADDRESS'", "'$ORACLE_ADDR'", "'$JOB_ID'");'
```

Hooray!

if we see failure, we can get transaction id and debug with truffle:
```bash
kubectl logs deploy/ganache
# get transaction id; go to the truffle directory containing the contract, and:
../../node_modules/.bin/truffle debug --network ganache <transaction id>
```


