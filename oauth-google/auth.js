const inquirer = require('inquirer');
const axios = require('axios').default;
const http = require('http')
const fs = require('fs')

let data = {
    client_id: "abc",
    client_secret: "def",
    refresh_token: "ghi"
}

fs.writeFile("test.json", JSON.stringify(data), function (err) {
    console.log(err)
})

// main()

async function main() {

    let base_url = "https://accounts.google.com/o/oauth2/v2/auth?scope=https://www.googleapis.com/auth/calendar https://www.googleapis.com/auth/calendar.events&access_type=offline&response_type=code&redirect_uri=http://localhost:4200&"

    let body = "redirect_uri=http://localhost:4200&grant_type=authorization_code"
    
    console.log("This will begin the process of the OAuth 2.0 Setup:")

    let question = [
        {
            name: "client_id",
            message: "Enter your OAuth Credential Client ID"
        },
        {
            name: "client_secret",
            message: "Enter your OAuth Credential Client Secret"
        }
    ]

    await inquirer.prompt(question).then(function (res) {
        base_url += "client_id=" + res.client_id
        body += "&client_id=" + res.client_id
        body += "&client_secret=" + res.client_secret
    })

    console.log("Go to the following URI\n" + base_url)

    let server = http.createServer(function (req,res) {
        let param = req.url.match(/code=([^&]*)/g)
        if (param && param.length > 0) {
            param = param[0].replace("code=","")
        }
        res.writeHead(200, { "Content-Type": "text/plain"}),
        res.write("Code=" + param),
        res.end()
    }).listen(4200)

    await inquirer.prompt({
        name: "code",
        message: "Please copy the code here"
    }).then(function (res) {
        body += "&code=" + res.code
        server.close()
    })

    console.log(encodeURI(body))

    axios.post("https://oauth2.googleapis.com/token", body).then(function (res) {
        console.log(res)
    }).catch(function (error) {
        console.log(error)
    })


}

