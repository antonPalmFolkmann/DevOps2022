export function getMessages() {
    var myHeaders = new Headers();
    myHeaders.append("Accept", "text/plain");

    var requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow'
    };

    return fetch("/public?offset=0&limit=200", requestOptions)
    .then(response => response.text())
    .then(result => JSON.parse(result))
    .then(result => exportMessages = result)
    .catch(error => console.log('error', error));

}

export function getMessagesByUser(user) {
    var myHeaders = new Headers();
    myHeaders.append("Accept", "text/plain");

    var requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow'
    };

    return fetch("/msgs/" + user, requestOptions)
    .then(response => response.text())
    .then(result => JSON.parse(result))
    .then(result => {
        console.log(result.Msgs)
        exportMessages = result.Msgs;
        exportProfile.username = result.username;
    })
    .catch(error => console.log('error', error));

}

export function postMessage(authorName, text) {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    console.log('In postmessage ' + authorName)
    console.log('In postmessage ' + text)


    var raw = JSON.stringify({
        "authorName": authorName,
        "text": text
    });

    var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
    };

    return fetch("/add_message", requestOptions)
        .then((response) => response.text())
        .then(result => {
            console.log('Result ' + result)
        })
        .catch(error => console.log('error', error));
}


export let exportProfile = {
    username: '',
    email: '',
    avatar: '',
    follows: [],
}
export let exportMessages = []