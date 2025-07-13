const http = require('http')

const num = Math.floor(Math.random() * 1000000)
const str = `string-${num}`

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' })
    const date = new Date().toJSON()
    res.end(`${date}: ${str}\n`)
})

const port = process.env.PORT || 3000
const host = '0.0.0.0'

server.listen(port, host, () => {
    console.log(`Server running at http://${host}:${port}/`)
})
