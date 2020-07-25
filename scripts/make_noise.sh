#!/bin/bash

geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "0xACa94ef8bD5ffEE41947b4585a84BdA5a3d3DA6E",to: "0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e", value: "10000000000000000000"})'
geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e",to: "0xACa94ef8bD5ffEE41947b4585a84BdA5a3d3DA6E", value: "10000000000000000000"})'
geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "0xACa94ef8bD5ffEE41947b4585a84BdA5a3d3DA6E",to: "0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e", value: "10000000000000000000"})'
geth attach http://localhost:32000 -exec 'eth.sendTransaction({from: "0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e",to: "0xACa94ef8bD5ffEE41947b4585a84BdA5a3d3DA6E", value: "10000000000000000000"})'