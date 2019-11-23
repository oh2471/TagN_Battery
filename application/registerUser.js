'use strict';

const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

class registerUser {

    main(data) {
        return new Promise(
            async (resolve, reject) => {
                try {
                    console.log('main data', data);
                    const walletPath = path.join(process.cwd(), 'wallet');
                    const wallet = new FileSystemWallet(walletPath);
                    console.log(`Wallet path: ${walletPath}`);
            
                    const adminExists = await wallet.exists('admin');
                    if (!adminExists) {
                        console.log('An identity for the admin user "admin" does not exist in the wallet');
                        console.log('Run the enrollAdmin.js application before retrying');
                        return;
                    }
            
                    const gateway = new Gateway();
                    await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: false } });
                    

                    const ca = gateway.getClient().getCertificateAuthority();
                    const adminIdentity = gateway.getCurrentIdentity();
                    
  
                    const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: data, role: 'client' }, adminIdentity);
                    const enrollment = await ca.enroll({ enrollmentID: data, enrollmentSecret: secret });
                    const userIdentity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
                    await wallet.import(data, userIdentity);
                    console.log(`Successfully registered and enrolled admin user ${data} and imported it into the wallet`);
                    resolve(true);
                } catch (error) {
                    reject(error);
                    console.error(`Failed to register user ${data}: ${error}`);
                    process.exit(1);
                }
            }
        )
    }
}

module.exports = new registerUser();