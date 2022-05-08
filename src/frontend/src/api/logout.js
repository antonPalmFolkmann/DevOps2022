export function logout() {

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
    });

    var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
    };

    return fetch("/logout", requestOptions)
        .then(response => JSON.stringify(response))
        .then(result => console.log(result))
        .catch(error => console.log('error', error));
}