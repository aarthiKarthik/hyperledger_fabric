#!/bin/bash

echo 
echo " "
echo "Starting marbles private data collection demo"
echo " "

CHANNEL_NAME="$1"
CC_NAME="$2"
VERSION_NUM="$3"
CC_SRC_PATH="github.com/chaincode/marbles02_private/go/"
COLLECTIONS_PATH="$GOPATH/src/github.com/chaincode/marbles02_private/collections_config.json"

echo "Channel name : "$CHANNEL_NAME

. scripts/utils.sh

installMarblesChaincode() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG
  set -x
  peer chaincode install -n ${CC_NAME} -v ${VERSION_NUM} -p ${CC_SRC_PATH} >&log.txt
  res=$?
  set +x
  cat log.txt
  verifyResult $res "Chaincode installation on peer${PEER}.org${ORG} has failed"
  echo "===================== Chaincode is installed on peer${PEER}.org${ORG} ===================== "
  echo
}

instantiateMarblesChaincode() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG

  # while 'peer chaincode' command can get the orderer endpoint from the peer
  # (if join was successful), let's supply it directly as we know it using
  # the "-o" option
  set -x
  peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_NAME} -v ${VERSION_NUM} -c '{"Args":["init"]}' -P "OR('Org1MSP.member','Org2MSP.member')" --collections-config  ${COLLECTIONS_PATH}
  res=$?
  set +x
  cat log.txt
  verifyResult $res "Chaincode instantiation on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' failed"
  echo "===================== Chaincode is instantiated on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' ===================== "
  echo
}

marbleschaincodeInvokeInit() {
  PEER=$1
  ORG=$2
  ARG=$3
  setGlobals $PEER $ORG
  PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

  set -x
  peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_NAME} -c '{"Args":["initMarble","'${ARG}'","blue","35","tom","99"]}'
  res=$?
  set +x
  cat log.txt
  verifyResult $res "Invoke execution on $PEERS failed"
  echo "===================== Invoke transaction successful on $PEERS on channel '$CHANNEL_NAME' ===================== "
  echo
}

## Install chaincode on peer0.org1 and peer0.org2
echo "Installing marbles chaincode on peer0.org1..."
installMarblesChaincode 0 1
echo "Install marbles chaincode on peer0.org2..."
installMarblesChaincode 0 2

# Instantiate chaincode on peer0.org1
echo "Instantiating chaincode on peer0.org1..."
instantiateMarblesChaincode 0 1

# Invoke chaincode on peer0.org1 and peer0.org2
echo "Sending invoke transaction on peer0.org1"
marbleschaincodeInvokeInit 0 1

