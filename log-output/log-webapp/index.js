const fs = require('fs')
const http = require('http')
const path = require('path')

const directory = path.join('/', 'usr', 'src', 'app', 'files')
const logPath = path.join(directory, 'log.txt')
const visitsPath = path.join(directory, 'visits.txt')

const server = http.createServer((req, res) => {
    if (fs.existsSync(logPath) && fs.existsSync(visitsPath)) {
        const log = fs.readFileSync(logPath)
        const visits = `Ping/pongs: ${fs.readFileSync(visitsPath)}`
        res.writeHead(200, { 'Content-Type': 'text/plain' })
        res.end(`${log}\n${visits}`)
    } else {
        res.writeHead(500, { 'Content-Type': 'text/plain' })
        res.end('Error: File not found')
    }
})

const port = process.env.PORT || 3000
const host = '0.0.0.0'

server.listen(port, host, () => {
    console.log(`Server running at http://${host}:${port}/`)
})
