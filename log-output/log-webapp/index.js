const fs = require('fs')
const http = require('http')
const path = require('path')

const directory = path.join('/', 'usr', 'src', 'app', 'files')
const filePath = path.join(directory, 'log.txt')

const server = http.createServer((req, res) => {
    
    if (fs.existsSync(filePath)) {
        const line = fs.readFileSync(filePath)
        res.writeHead(200, { 'Content-Type': 'text/plain' })
        res.end(line)
    } else {
        res.writeHead(500, { 'Content-Type': 'text/plain' })
        res.end("Error: File not found")
    }
})

const port = process.env.PORT || 3000
const host = '0.0.0.0'

server.listen(port, host, () => {
    console.log(`Server running at http://${host}:${port}/`)
})
