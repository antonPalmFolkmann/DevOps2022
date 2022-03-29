const scanLicenseResults = require('./license.json')
const scanSourceResults  = require('./src.json')
let nrOfLicensesFound = 0

console.log('Checking ScanCode result for licenses')

for (let index of scanLicenseResults.files.keys()) {
    if (scanLicenseResults.files[index].licenses.length > 0) {
        console.log('License(s) found in file: ' + JSON.stringify(scanLicenseResults.files[index].path))
        for (let jndex of scanLicenseResults.files[index].licenses.keys()) {
            nrOfLicensesFound++
            console.log(scanLicenseResults.files[index].licenses[jndex].name)
        }
    }
}

for (let index of scanSourceResults.files.keys()) {
    if (scanSourceResults.files[index].licenses.length > 0) {
        console.log('License(s) found in file: ' + JSON.stringify(scanSourceResults.files[index].path))
        for (let jndex of scanSourceResults.files[index].licenses.keys()) {
            nrOfLicensesFound++
            console.log(scanSourceResults.files[index].licenses[jndex].name)
        }
    }
}

if (nrOfLicensesFound > 0) {
    console.log('Check complete. ' + nrOfLicensesFound + ' licenses found.')
} else {
    console.log('Check complete. No licenses found.')
}