const http = require('http')
const port = 3000

const requestHandler = (request, response) => {
    var data = ''
    request.on('data', chunk => {
        data += chunk
    });
    request.on('end', () => {
        d = JSON.stringify({
            result: {
              availableVerifications: [{
                verificationType: "passport",
                verificationStatusType: "verified",
              }],
            },
        })
        console.log(request.method, request.url, data, d)
        response.end(d)
    })
}

const server = http.createServer(requestHandler)

server.listen(port, (err) => {
    if (err) {
        return console.log('something bad happened', err)
    }
    console.log(`server is listening on ${port}`)
})
