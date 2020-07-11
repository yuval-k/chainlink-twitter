npm install --save ganache-cli web3 solc

go to chainlink source and do `make install` to build the contracts



helm template release bitnami/postgresql > manifests/postgresql.yaml


you will see:
  LinkToken: 0x5b1869d9a4c187f2eaa108f3062412ecf0526b24

account owner


```bash
LINK_TOKEN=0x5b1869d9a4c187f2eaa108f3062412ecf0526b24
ADDRESS=0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1
KEY=0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d

make deploy-testnet
kubectl rollout status deploy/ganache
make deploy-token
# verify that Link token 0x5b1869d9a4c187f2eaa108f3062412ecf0526b24
make deploy-node

#wait for node:
kubectl rollout status deploy/chainlink


export NODE_ADDR=$(kubectl logs deploy/chainlink|grep "please deposit ETH into your address:"| tr ' ' '\n'|grep 0x)

geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "'${ADDRESS}'",to: "'${NODE_ADDR}'", value: "74000000000000000"})'


geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "'${ADDRESS}'",to: "'${NODE_ADDR}'", value: "74000000000000000"})'

# if link token is not 0x5b1869d9a4c187f2eaa108f3062412ecf0526b24
# update it in the oracle/migrations/2_deploy_contracts.js
npm run deploy-oracle | tee node-tmp.txt
export ORACLE_ADDR=$(grep "contract-address" node-tmp.txt | cut -f 2)
rm node-tmp.txt

kubectl port-forward deploy/chainlink 6688&
```

this part is silly.
log-in in the UI and then open a js console and paste:

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

paste it in the terminal before continuing

```bash
curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"ethbytes32"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httppost"},{"type":"jsonparse"},{"type":"ethbytes32"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"multiply"},{"type":"ethint256"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"multiply"},{"type":"ethuint256"},{"type":"ethtx"}]}'

curl http://localhost:6688/v2/specs -XPOST -H"X-API-KEY: $ACCESS_KEY" -H"X-API-SECRET: $SECRET_KEY" -H"content-type: application/json" -d '{"initiators":[{"type":"runlog","params":{"address":"'$ORACLE_ADDR'"}}],"tasks":[{"type":"httpget"},{"type":"jsonparse"},{"type":"ethbool"},{"type":"ethtx"}]}'


# fund the node:
geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("fund.js");transfer("'$LINK_TOKEN'", "'$ADDRESS'", "'$NODE_ADDR'");'


npm run deploy-testconsumer | tee node-tmp.txt
export TEST_CONTRACT_ADDR=$(grep "contract-address" node-tmp.txt | cut -f 2)
rm node-tmp.txt
# fund the contact:
geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("fund.js");transfer("'$LINK_TOKEN'", "'$ADDRESS'", "'$TEST_CONTRACT_ADDR'");'

# job id of EthUint256
# TODO: grab that from curl
JOB_ID="0xcb69fe2b09b344a581efc5338fd5f8b2"

# create a script with our test contract
echo "var contract_ = " $(cat contracts/testconsumer/build/contracts/ATestnetConsumer.json|jq .abi) ";function contract() {return contract_;} " > scripts/contract.js

geth attach http://localhost:32000 --jspath ./scripts -exec 'loadScript("contract.js");loadScript("request.js");request(contract(), "'$TEST_CONTRACT_ADDR'", "'$ORACLE_ADDR'", "'$JOB_ID'");'


```


kubectl port-forward deploy/chainlink 6688

go to http://localhost:6688
yes!


move some eth to the node using metamask
node address is in config page
TODO: can i automate this: move eth to node using command line?
move link to node using command line?
get node address via command line

deploy oracle and call setFulfillmentPermission via script / commandline




now move some link or ether to node:

do this https://docs.chain.link/docs/fulfilling-requests

deploy oracle:
add local test net to metamask
add link token ot moraclee plugins
injected web 3
at addess = link token above