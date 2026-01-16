import * as fs from 'node:fs/promises'
import http from 'http'
import path from 'path'

const logPath = path.join('/', 'usr', 'src', 'app', 'files', 'log.txt')
const infoPath = path.join('/', 'information', 'information.txt')
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
    if (req.url === '/ready') {
        try {
            const data = await getRequest(`${PING_PONG_SVC_URL}/ping-count`)
            res.writeHead(200, { 'Content-Type': 'application/json' })
            res.end(JSON.stringify({ status: 'ready' }))
        } catch {
            const msg = 'Error: Failed to retrieve ping-pong count'
            res.writeHead(503, { 'Content-Type': 'application/json' })
            res.end(JSON.stringify({ error: 'Ping pong service not ready' }))
        }
        return
    }

    let log = ''
    let info = ''
    try {
        log = await fs.readFile(logPath)
        info = await fs.readFile(infoPath)
    } catch {
        const msg = 'Error: File not found'
        console.log(msg)
        res.writeHead(500, { 'Content-Type': 'text/plain' })
        res.end(msg)
        return
    }

    try {
        const data = await getRequest(`${PING_PONG_SVC_URL}/ping-count`)
        res.writeHead(200, { 'Content-Type': 'text/plain' })
        const responseData = [
            'File content:',
            info,
            `MESSAGE=${process.env.MESSAGE}`,
            log,
            `Ping/pongs: ${data.pings}`
        ]
        res.end(responseData.join('\n'))
    } catch {
        const msg = 'Error: Failed to retrieve ping-pong count'
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
