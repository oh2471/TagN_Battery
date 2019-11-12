
docker exec cli peer chaincode install -n elca -v 1.0 -p github.com/test/
export CORE_PEER_ADDRESS=peer1.org1.battery.com:7051
docker exec cli peer chaincode install -n elca -v 1.0 -p github.com/test/

export CORE_PEER_LOCALMSPID=Org2MSP
export PEER0_ORG2_CA=/etc/hyperledger/crypto/peerOrganizations/org2.battery.com/peers/peer0.org2.battery.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peerOrganizations/org2.battery.com/users/Admin@org2.battery.com/msp
export CORE_PEER_ADDRESS=peer0.org2.battery.com:7051
docker exec cli peer chaincode install -n elca -v 1.0 -p github.com/test/
export CORE_PEER_ADDRESS=peer1.org2.battery.com:7051
docker exec cli peer chaincode install -n elca -v 1.0 -p github.com/test/


export CORE_PEER_LOCALMSPID=Org1MSP
export PEER0_ORG1_CA=/etc/hyperledger/crypto/peerOrganizations/org1.battery.com/peers/peer0.org1.battery.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peerOrganizations/org1.battery.com/users/Admin@org1.battery.com/msp
export CORE_PEER_ADDRESS=peer0.org1.battery.com:7051

export ORDERER_CA=/etc/hyperledger/crypto/ordererOrganizations/battery.com/orderers/orderer.battery.com/msp/tlscacerts/tlsca.battery.com-cert.pem
#docker exec cli peer chaincode instantiate -o orderer.battery.com:7050 -C battery -n elca -v 1.0 -c '{"Args":["init"]}' -P "OR('Org1MSP.member','Org2MSP.member')" --collections-config /opt/gopath/src/github.com/battery/collections_config.json


#--tls --cafile $ORDERER_CA 
# marbles
# export MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)
# export MARBLE=$(echo -n "{\"name\":\"marble2\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)
# docker exec cli peer chaincode invoke -o orderer.battery.com:7050 -C battery -n elca -c '{"Args":["initMarble"]}'  --transient "{\"marble\":\"$MARBLE\"}"

# docker exec cli peer chaincode query -C battery -n elca -c '{"Args":["readMarble","marble1"]}'
# docker exec cli peer chaincode query -C battery -n elca -c '{"Args":["readMarblePrivateDetails","marble1"]}'


# Chaincode example (private data)
# docker exec cli peer chaincode invoke -n elca -C battery -c '{"Args":["saveinfo","1","2"]}'

# export PI=$(echo -n "{\"gps\":\"123.412.312\"}" | base64 | tr -d \\n)
# docker exec cli peer chaincode invoke -o orderer.battery.com:7050 -C battery -n elca -c '{"Args":["savePersonalInfo"]}'  --transient "{\"personalInfo\":\"$PI\"}"
# docker exec cli peer chaincode query -C battery -n elca -c '{"Args":["getPersonalInfo","LEE"]}'

docker exec cli peer chaincode instantiate -v 1.0 -C battery -n elca -c '{"Args":["a","100"]}' -P 'OR ("Org1MSP.member", "Org2MSP.member")'
docker exec cli peer chaincode invoke -o orderer.battery.com:7050 -C battery -n elca -c '{"Args":["set","b","200"]}'
docker exec cli peer chaincode query -C battery -n elca -c '{"Args":["get","b"]}'




#docker logs peer0.org2.battery.com 2>&1 | grep -i -a -E 'private|pvt|privdata'

echo '-------------------------------------END-------------------------------------'

