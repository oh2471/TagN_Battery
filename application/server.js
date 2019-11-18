const express = require('express');
const app = express();
const path = require('path');
const bodyParser = require('body-parser');
const main = require('./router/main');
const data = require('./router/data');
const static = require('serve-static');
const session = require('express-session');
const FileStore = require('session-file-store')(session);

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

