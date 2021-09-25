const http = require('http');
const util = require('util');
const fs = require('fs');

const port = 3000;

let config = [{
  urlRe: /\/search/,
  respFile: 'serp.json',
}];

const logJSON = (v) => {
  console.log(util.inspect(JSON.parse(v), {
    depth: null,
    colors: true,
    breakLength: 160,
    sorted: true,
  }));
};

const requestHandler = (request, response) => {
  const { headers, method, url } = request;
  var data = '';
  request.on('data', chunk => {data += chunk});
  request.on('end', () => {
    var resp = undefined;
    for(var i = 0; i < config.length; i++) {
      var c = config[i];
      if (c.urlRe.test(url)) {
        console.log('matched', c);
        resp = fs.readFileSync(c.respFile);
        break;
      } else {
        console.log('skipped', c);
      }
    }
    console.log(`\x1b[1;44;36m${url}\x1b[K\x1b[0m`);
    logJSON(data);
    logJSON(resp);
    response.end(resp);
  });
}

const server = http.createServer(requestHandler);

server.listen(port, (err) => {
    if (err) {
        return console.log('something bad happened', err);
    }
    console.log(`server is listening on ${port}`);
});
