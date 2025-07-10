const num = Math.floor(Math.random() * 1000000)
const str = `string-${num}`

setInterval(() => {
    const date = new Date().toJSON();
    console.log(`${date}: ${str}`)
}, 5000)
