module.exports = {
  networks: {
    ganache: {
      host: "127.0.0.1",
      port: 32000,
      network_id: "*",
    },
    development: {
      host: "127.0.0.1",
      port: 32000,
      network_id: "*"
    },
    test: {
      host: "127.0.0.1",
      port: 32000,
      network_id: "*"
    }
  },
  compilers: {
    solc: {
      version: "0.4.24",
    }
  }  
};
