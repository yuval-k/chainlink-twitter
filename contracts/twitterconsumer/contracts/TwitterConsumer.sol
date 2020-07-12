pragma solidity 0.4.24;

import "../../vendor/v0.4/ChainlinkClient.sol";
import "../../vendor/v0.4/vendor/Ownable.sol";


/**
The goal of this contract is to allow allow the transfer of ETH from the originator to the
beneficiary for real world sevices. To make sure all goes well, the funds will only transfer after
a trusted 3rd party approves it.
The trusted 3rd party is identified by his twitter handle. and he approves the transaction by 
tweeting the approval text.
A custom chainlink adapter checks twitter and confirms the contract.

The flow is as follows:
1. The originator creates the contract. Both parties can see terms of the contract: 
   - The identity (twitter handle) of the approver
   - The amount.
   - The orcale / job used to verify
   - The expiration date
1. In the real world the beneficiary confirms satisfaction with the contract.
1. The originator funds the contract with ETH (calls `fund()` with the exact amount agreed upon).
1. The originator funds the contract with LINK.
1. The beneficiary can check that the contract is ready using `ready()`.
1. Now that the contract is ready, the real world transaction can happen
1. Trusted 3rd party tweets the approval text.
1. Someone (originator or beneficiary) calls `requestApproval()`
1. Chainlink magic happens, during which the node uses twitter API to see if the approval text as tweeted.
1. Assuming the transaction as approved, the beneficiary can now call withdraw to receive the funds.
1. If the contract expires and no approval was given, the originator can refund the contract.
 */
contract TwitterConsumer is ChainlinkClient, Ownable {
    uint256 private constant ORACLE_PAYMENT = 1 * LINK;

    address payable public originator;
    address payable public beneficiary;
    uint256 public deadline;
    uint256 public amount;

    string public handle;
    string public text;
    address public oracle;
    string public jobId;

    // done means that the node returned an answer, and the contract reached it'se final state.
    bool done;
    // approved means that the trusted 3rd party approved the transaction.
    bool approved;

    modifier onlyBeneficiary() {
        require(msg.sender == beneficiary);
        _;
    }

    modifier onlyOriginator() {
        require(msg.sender == originator);
        _;
    }
    event Funded(uint256 amount);
    event Fulfilled(bytes32 _requestId, bool allow);

    constructor(
        address _link,
        uint256 _deadline,
        address payable _beneficiary,
        uint256 amount,
        string _approver_twitter_handle,
        string _text,
        address _oracle,
        string _jobId
    ) public {
        setChainlinkToken(_link);
        originator = msg.sender;
        beneficiary = _beneficiary;
        amount = _amount;
        deadline = now + _deadline;
        handle = _approver_twitter_handle;
        text = _text;
        oracle = _oracle;
        jobId = _jobId;
    }

    // this will be called when ETH is sent to the contract automatically.
    receive() external payable {
        fund();
    }

    function fund() public onlyOriginator payable {
        require(!done, "Contract is done");
        require(amount == 0, "Contract already fudned");
        require(msg.value == amount, "Not the amount agreed upon");
        require(now <= deadline, "Deadline expired");
        emit Funded(msg.value);
    }

    function ready() public view returns (bool) {
        return
            (link.balanceOf(address(this)) >= ORACLE_PAYMENT) && (balance == amount);
    }

    function requestApproval()
        public
    {
        require(!done, "Contract is done");
        Chainlink.Request memory req = buildChainlinkRequest(
            stringToBytes32(_jobId),
            this,
            this.fulfillApproval.selector
        );
        req.add("handle", handle);
        req.add("text", text);
        sendChainlinkRequestTo(_oracle, req, ORACLE_PAYMENT);
    }

    function fulfillApproval(bytes32 _requestId, bool _done, bool _approved)
        public
        recordChainlinkFulfillment(_requestId)
    {
        if (!done) {
          return;
        }
        done = true;
        approved = _approved;
        emit Fulfilled(_requestId, _approved);
    }

    function withdraw() public onlyBeneficiary {
        require(done, "Contract is still in progress");
        require(approved == true, "Cannot withdraw an unapproved contract");
        selfdestruct(beneficiary);
    }

    function refund() public onlyOriginator {
        if (!done && now > deadline) {
            done = true;
            // we have reached deadline, allow refund.
        } else {
            require(done, "Contract is still in progress");
            require(approved == false, "Cannot refund an approved contract");
        }
        selfdestruct(originator);
    }

    function getChainlinkToken() public view returns (address) {
        return chainlinkTokenAddress();
    }

    function withdrawLink() public onlyOriginator {
        LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
        require(
            link.transfer(msg.sender, link.balanceOf(address(this))),
            "Unable to transfer"
        );
    }

    function stringToBytes32(string memory source)
        private
        pure
        returns (bytes32 result)
    {
        bytes memory tempEmptyStringTest = bytes(source);
        if (tempEmptyStringTest.length == 0) {
            return 0x0;
        }

        assembly { // solhint-disable-line no-inline-assembly
            result := mload(add(source, 32))
        }
    }
}
