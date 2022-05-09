export function login(username, password) {

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
        "username": username,
        "password": password
    });

    var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
    };

    return fetch("/login", requestOptions)
        .then((response) => response.text())
        .then(result => {
            if (result[0] === '{') {
                loggedInUser = JSON.parse(result)
            } else {
                loggedInUser = {
                    username: '',
                }
            }
        })
        .catch(error => console.log('error', error));
}

export let loggedInUser = {}