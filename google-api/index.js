const { google } = require('googleapis')

runTest()

async function runTest() {
    // const auth = new google.auth.GoogleAuth({
    //     keyFile: path.join("Key File Path"),
    //     scopes: ["https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/calendar.events"]
    // })

    const jwt = new google.auth.JWT({
        email: "Email",
        key: "Private Key",
        scopes: ["https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/calendar.events"]
    })

    // const client = await auth.getClient();
    // jwt.authorize(function (err, res) {
    //     let events = google.calendar({
    //         version: 'v3',
    //         auth: jwt
    //     }).events.list({
    //         calendarId: "email address",
    //         timeMin: "2020-10-28T07:59:00-07:00",
    //         timeMax: "2020-10-28T15:00:00-07:00"
    //     })

    //     console.log(events)
    // })
    const calendar = google.calendar({
        version: 'v3',
        auth: jwt,
        
    })

    const res = await calendar.events.list({
        calendarId: "email address",
        timeMin: "2020-10-28T07:59:00-07:00",
        timeMax: "2020-10-28T15:00:00-07:00"
    })

    let availableTimes = [8,9,10,11,12,13,14,15]
    console.log(__dirname)
    if (res.data && res.data.items) {
        let events = res.data.items
        events.forEach(function(event) {
            let date = new Date(event.start.dateTime)
            let hour = date.getHours()
            let index = availableTimes.indexOf(hour)
            if (index > -1) {
                availableTimes.splice(index, 1)
            }
        })
    }

    console.log(availableTimes)
    
    // const res = await calendar.events.insert({
    //     calendarId: "email address",
    //     requestBody: {
    //         summary: "This is my test event",
    //         start: {
    //             dateTime: "2020-10-29T10:00:00-08:00",
    //             timeZone: "America/Los_Angeles"
    //         },
    //         end: {
    //             dateTime: "2020-10-29T11:00:00-08:00",
    //             timeZone: "America/Los_Angeles"
    //         }
    //     }
    // })

    // console.log(res.data);

    return res.data
}
 