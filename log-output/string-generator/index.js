const fs = require('fs')
const path = require('path')

const num = Math.floor(Math.random() * 1000000)
const str = `string-${num}`

const directory = path.join('/', 'usr', 'src', 'app', 'files')
const filePath = path.join(directory, 'log.txt')

if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory)
}

setInterval(() => {
    const date = new Date().toJSON()
    const line = `${date}: ${str}`
    console.log(line)
    fs.writeFileSync(filePath, line)
}, 5000)
