docker exec cli peer chaincode install -n elca -v 1.0 -p github.com/battery/
export CORE_PEER_LOCALMSPID=Org2MSP
export PEER0_ORG2_CA=/etc/hyperledger/crypto/peerOrganizations/org2.battery.com/peers/peer0.org2.battery.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peerOrganizations/org2.battery.com/users/Admin@org2.battery.com/msp
export CORE_PEER_ADDRESS=peer0.org2.battery.com:7051
docker exec cli peer chaincode install -n elca -v 1.0 -p github.com/battery/

#export ORDERER_CA=/etc/hyperledger/crypto/ordererOrganizations/battery.com/orderers/orderer.battery.com/msp/tlscacerts/tlsca.battery.com-cert.pem
docker exec cli peer chaincode instantiate -o orderer.battery.com:7050 -C battery -n elca -v 1.0 -c '{"Args":["init"]}' -P "OR('Org1MSP.member','Org2MSP.member')" --collections-config /opt/gopath/src/github.com/battery/collections_config.json

#--tls --cafile $ORDERER_CA 

export MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)

export MARBLE=$(echo -n "{\"name\":\"marble2\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)
docker exec cli peer chaincode invoke -o orderer.battery.com:7050 -C battery -n elca -c '{"Args":["initMarble"]}'  --transient "{\"marble\":\"$MARBLE\"}"

docker exec cli peer chaincode query -C battery -n elca -c '{"Args":["readMarble","marble1"]}'
docker exec cli peer chaincode query -C battery -n elca -c '{"Args":["readMarblePrivateDetails","marble1"]}'



#docker logs peer0.org1.battery.com 2>&1 | grep -i -a -E 'private|pvt|privdata'

echo '-------------------------------------END-------------------------------------'

