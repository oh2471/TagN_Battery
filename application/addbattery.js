'use strict';

const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const registerUser = require("./registerUser");
const ccpPath = path.resolve(__dirname, '..', 'network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

class addbattery {

    main(data) {
        return new Promise(
            async (resolve, reject) => {
                try {
                    
                    const walletPath = path.join(process.cwd(), 'wallet');
                    const wallet = new FileSystemWallet(walletPath);
                    console.log(`Wallet path: ${walletPath}`);

                    const userExists = await wallet.exists(data.phone);
                    if (!userExists) {
                        await registerUser.main(data.phone)
                    }
                    const gateway = new Gateway();
                    await gateway.connect(ccp, { wallet, identity: data.phone, discovery: { enabled: false } });
    
                    const network =await gateway.getNetwork('battery');
                    const contract = network.getContract('elca');



                    await contract.submitTransaction(
                        'addBattery',
                        data.phone.toString(),
                        data.batterystatus.toString(),
                        data.stationname.toString(),
                        data.stationgps.toString(),
                        data.check.toString()
                    );

                    console.log("submitTransaction Success");
                    await gateway.disconnect();
                    resolve(true);
                } catch (error) {
                    reject(error);
                    process.exit(1);
                }
            }
        )
    }
}

module.exports = new addbattery();