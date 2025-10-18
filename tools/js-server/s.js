const http = require('http');
const util = require('util');
const fs = require('fs');

const port = process.argv[2] || 3000;

function logJSON(v) {
  try {
    console.log(util.inspect(JSON.parse(v), {
      depth: null,
      colors: true,
      breakLength: 160,
      sorted: true,
    }));
  } catch (err) {
    console.log(v);
  }
}

function readFile(filename) {
  return new Promise((resolve, reject) => {
    fs.readFile(filename, 'utf8', (err, data) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(data);
    });
  });
}

function proxyCall(options, data) {
  logJSON(options);
  return new Promise((resolve, reject) => {
    const req = http.request(options, (resp) => {
      let d = '';
      resp.on('data', (chunk) => { d += chunk; });
      resp.on('end', () => { resolve(d); });
      resp.on('error', (err) => { reject(err); });
    });
    req.write(data);
    req.end();
  });
}

function requestHandler(request, response) {
  const { headers, method, url } = request;
  console.log(`\x1b[1;44;30m${method} \x1b[96m${url}\x1b[30m ${new Date().toISOString().replace(/T/, ' ').replace(/Z/, '')}\x1b[K\x1b[0m`);
  let data = '';
  request.on('data', (chunk) => { data += chunk; });
  request.on('end', async () => {
    let config = JSON.parse(await readFile('config.json'));
    let resp;
    let code = 200;
    let headers = {};
    for (let i = 0; i < config.length; i++) { // TODO no-await-in-loop
      const c = config[i];
      if ((new RegExp(c.urlRe)).test(url)) {
        console.log('Matched', c.urlRe);
        if (c.respFile) {
          resp = await readFile(c.respFile);
        } else if (c.payload) {
          resp = JSON.stringify(c.payload)
        } else if (c.location) {
          code = 302;
          headers['Location'] = c.location;
          resp = '';
        } else if (c.proxyHost) {
          const h = { ...headers };
          h.host = c.proxyHost;
          const opts = {
            host: c.proxyHost,
            headers: h,
            method,
            path: url,
            port: c.proxyPort || 80,
            setHost: true,
          };
          resp = await proxyCall(opts, data);
        }
        if (c.mime) {
          headers['Content-Type'] = c.mime;
        }
        break;
      } else {
        console.log('Skipped', c.urlRe);
      }
    }
    logJSON(data);
    logJSON(resp);
    response.writeHead(code, headers);
    response.end(resp);
  });
}

const server = http.createServer(requestHandler);

server.listen(port, (err) => {
  if (err) {
    console.log('Something bad happened', err);
  }
  console.log(`Server is listening on ${port}`);
});
