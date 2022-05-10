export function register(username, email, password) {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    console.log('Username in register ' + username)
    console.log('Email in register ' + email)
    console.log('Password in register ' + password)



    var raw = JSON.stringify({
        "username": username,
        "email": email,
        "password": password
    });

    var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
    };

    fetch("/api/register", requestOptions)
        .then(response => JSON.stringify(response))
        .then(result => console.log(result))
        .catch(error => console.log('error', error));
}