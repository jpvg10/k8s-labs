import { stat, mkdir, writeFile } from 'node:fs/promises'
import { join } from 'node:path'

const num = Math.floor(Math.random() * 1000000)
const str = `string-${num}`

const directory = join('/', 'usr', 'src', 'app', 'files')
const filePath = join(directory, 'log.txt')

try {
    await stat(directory)
} catch (error) {
    await mkdir(directory)
}

setInterval(async () => {
    const date = new Date().toJSON()
    const line = `${date}: ${str}`
    console.log(line)
    await writeFile(filePath, line)
}, 5000)
