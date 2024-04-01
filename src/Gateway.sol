// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import "suave-std/suavelib/Suave.sol";
import "suave-std/Transactions.sol";
import "suave-std/protocols/Bundle.sol";
import "solady/src/utils/JSONParserLib.sol";
import "forge-std/console.sol";

contract Gateway {
    event NewBuilderBidEvent(Suave.DataId dataId, uint64 decryptionCondition, address[] allowedPeekers, bytes envelope);

    string boostRelayUrl;

    constructor(string memory boostRelayUrl_) {
        boostRelayUrl = boostRelayUrl_;
    }

    function emitNewBuilderBidEvent(Suave.DataRecord memory record, bytes memory envelope) public {
        emit NewBuilderBidEvent(record.id, record.decryptionCondition, record.allowedPeekers, envelope);
    }

    function buildAndSendToRelay(
        uint64 blockNumber,
        Suave.BuildBlockArgs memory blockArgs,
        address[] calldata allowedPeekers,
        address[] calldata allowedStores
    ) external payable returns (bytes memory) {
        require(Suave.isConfidential());
        Suave.DataId[] memory bundleDataIds = new Suave.DataId[](0); // empty array - will be filled by pending tx on builder.
        Suave.DataRecord memory record =
            Suave.newDataRecord(blockNumber, allowedPeekers, allowedStores, "default:v0:ethBlocks");
        Suave.confidentialStore(record.id, "default:v0:mergedDataRecords", abi.encode(bundleDataIds));
        (bytes memory builderBid, bytes memory envelope) = Suave.buildEthBlock(blockArgs, record.id, ""); // namespace not used.
        Suave.confidentialStore(record.id, "default:v0:builderBids", builderBid);
        Suave.confidentialStore(record.id, "default:v0:payloads", envelope);
        Suave.submitEthBlockToRelay(boostRelayUrl, builderBid);
        return bytes.concat(this.emitNewBuilderBidEvent.selector, abi.encode(record, envelope));
    }
}
