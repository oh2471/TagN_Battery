const express = require('express');
const app = express();
const path = require('path');
const bodyParser = require('body-parser');
const main = require('./router/main');
const data = require('./router/data');
const static = require('serve-static');
const session = require('express-session');
const FileStore = require('session-file-store')(session);

// 하이퍼레저 연동
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);


app.use(session({
    secret: 'keyboard cat',
    resave: false,
    saveUninitialized: true,
    store: new FileStore(),
}));

app.set("/views", static(path.join(__dirname, 'views')));

app.listen(8080, function() {
    console.log("start! server on port 8080");
});

app.use(express.static('public'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended:true}));
app.set('view engine', 'ejs');
app.engine('html', require('ejs').renderFile);
app.use('/main',main);
app.use(main);
app.use('/data',data);


async function cc_call(fn_name, args){
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('teamate');

    var result;
    
    if(fn_name == 'addUser')
        result = await contract.submitTransaction('addUser', args);
    else if( fn_name == 'addRating')
    {
        e=args[0];
        p=args[1];
        s=args[2];
        result = await contract.submitTransaction('addRating', e, p, s);
    }
    else if(fn_name == 'readRating')
        result = await contract.evaluateTransaction('readRating', args);
    else
        result = 'not supported function'

    return result;
}
