import * as fs from 'node:fs/promises'
import http from 'http'
import path from 'path'

const directory = path.join('/', 'usr', 'src', 'app', 'files')
const logPath = path.join(directory, 'log.txt')
const { 
    PING_PONG_SVC_URL = 'http://localhost:4000',
    PORT = '3000'
} = process.env

const getRequest = (url) => {
    return new Promise((resolve, reject) => {
        const req = http.get(url, (res) => {
            if (res.statusCode !== 200) {
                reject()
            }

            let data = ''
            res.on('data', (chunk) => {
                data += chunk
            })

            res.on('close', () => {
                resolve(JSON.parse(data))
            })
        })

        req.on('error', (err) => {
            reject(err)
        })
    })
}

const server = http.createServer(async (req, res) => {
    try {
        const log = await fs.readFile(logPath)

        try {
            const data = await getRequest(`${PING_PONG_SVC_URL}/ping-count`)
            res.writeHead(200, { 'Content-Type': 'text/plain' })
            res.end(`${log}\nPing/pongs: ${data.pings}`)
        } catch {
            const msg = 'Error: Failed to retrieve ping-pong count'
            console.log(msg)
            res.writeHead(500, { 'Content-Type': 'text/plain' })
            res.end(msg)
        }
    } catch {
        const msg = 'Error: File not found'
        console.log(msg)
        res.writeHead(500, { 'Content-Type': 'text/plain' })
        res.end(msg)
    }
})

const host = '0.0.0.0'
server.listen(PORT, host, () => {
    console.log(`Server running at http://${host}:${PORT}/`)
    console.log(`Ping-pong service URL: ${PING_PONG_SVC_URL}`)
})
