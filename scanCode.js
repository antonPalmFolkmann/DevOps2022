const scanLicenseResults = require('./license.json')
const scanSourceResults  = require('./src.json')
const scanUrlsResults = require('./urls.json')
let nrOfLicensesFound = 0
let nrOfUrlsFound = 0
let licensesFound = new Array()
let urlsFound = new Array()

console.log('CHECKING SCANCODE RESULT FOR LICENSES:')

for (let index of scanLicenseResults.files.keys()) {
    if (scanLicenseResults.files[index].licenses.length > 0) {
        console.log('License(s) found in file: ' + JSON.stringify(scanLicenseResults.files[index].path))
        for (let jndex of scanLicenseResults.files[index].licenses.keys()) {
            nrOfLicensesFound++
            console.log(scanLicenseResults.files[index].licenses[jndex].name)
            licensesFound.push(scanLicenseResults.files[index].licenses[jndex].name)
        }
    }
}

console.log(' \n')

for (let index of scanSourceResults.files.keys()) {
    if (scanSourceResults.files[index].licenses.length > 0) {
        console.log('License(s) found in file: ' + JSON.stringify(scanSourceResults.files[index].path))
        console.log(' \n')
        for (let jndex of scanSourceResults.files[index].licenses.keys()) {
            nrOfLicensesFound++
            console.log(scanSourceResults.files[index].licenses[jndex].name)
            licensesFound.push(scanSourceResults.files[index].licenses[jndex].name)
        }
        console.log(' \n')
    }
}


console.log('CHECKING SCANCODE RESULT FOR URLS:')

for (let index of scanUrlsResults.files.keys()) {
    if (scanUrlsResults.files[index].urls.length > 0) {
        for (let jndex of scanUrlsResults.files[index].urls.keys()) {
            nrOfUrlsFound++
            console.log(scanUrlsResults.files[index].urls[jndex].url)
            urlsFound.push(scanUrlsResults.files[index].urls[jndex].url)
        }
    }
}

console.log(' \n')

for (let index of licensesFound.keys()) {
    if (licensesFound[index] != 'MIT License') {
        throw new Error('License other than MIT found. Please review terminal output.')
    }
}

if (nrOfLicensesFound > 0) {
    console.log('Check complete. ' + nrOfLicensesFound + ' licenses found.')
} else {
    console.log('Check complete. No licenses found.')
}

if (nrOfUrlsFound > 0) {
    console.log('' + nrOfUrlsFound + ' urls found.')
} else {
    console.log('No urls found.')
}