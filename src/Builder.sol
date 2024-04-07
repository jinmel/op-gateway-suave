// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import "suave-std/suavelib/Suave.sol";

contract Builder {
    event NewBuilderBidEvent(Suave.DataId dataId, uint64 decryptionCondition, address[] allowedPeekers, bytes envelope);

    event NewEvent(uint64 n);

    function emitNewBuilderBidEvent(Suave.DataRecord memory record, bytes memory envelope) public {
        emit NewBuilderBidEvent(record.id, record.decryptionCondition, record.allowedPeekers, envelope);
    }

    function callback(uint n) public {
        emit NewEvent(uint64(n));
    }

    function build(
        uint64 blockNumber,
        string calldata relayUrl,
        address[] calldata allowedPeekers,
        address[] calldata allowedStores
    ) external payable returns (bytes memory) {
        require(Suave.isConfidential());
        Suave.DataId[] memory dataids = new Suave.DataId[](0);
        Suave.DataRecord memory record =
            Suave.newDataRecord(blockNumber, allowedPeekers, allowedStores, "random");
        Suave.confidentialStore(record.id, "default:v0:mergedDataRecords", abi.encode(dataids));
        Suave.BuildBlockArgs memory blockArgs;
        blockArgs.fillPending = true;
        (bytes memory builderBid, bytes memory envelope) = Suave.buildEthBlock(blockArgs, record.id, ""); // namespace not used.
        Suave.confidentialStore(record.id, "default:v0:builderBids", builderBid);
        Suave.confidentialStore(record.id, "default:v0:payloads", envelope);
        Suave.submitEthBlockToRelay(relayUrl, builderBid);
        return bytes.concat(this.emitNewBuilderBidEvent.selector, abi.encode(record, envelope));
    }
}
