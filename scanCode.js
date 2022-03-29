const scanLicenseResults = require('./license.json')
const scanSourceResults  = require('./src.json')
let nrOfLicensesFound = 0

console.log('Checking ScanCode result for licenses')

for (let i = 0; i < scanLicenseResults.files.length; i++) {
    if (scanLicenseResults.files[i].licenses.length > 0) {
        console.log('License(s) found in file: ' + JSON.stringify(scanLicenseResults.files[i].path))
        for (let j = 0; j < scanLicenseResults.files[i].licenses.length; j++) {
            nrOfLicensesFound++
            console.log(scanLicenseResults.files[i].licenses[j].name)
        }
    }
}

for (let i = 0; i < scanSourceResults.files.length; i++) {
    if (scanSourceResults.files[i].licenses.length > 0) {
        console.log('License(s) found in file: ' + JSON.stringify(scanSourceResults.files[i].path))
        for (let j = 0; j < scanSourceResults.files[i].licenses.length; j++) {
            nrOfLicensesFound++
            console.log(scanSourceResults.files[i].licenses[j].name)
        }
    }
}

if (nrOfLicensesFound > 0) {
    console.log('Check complete. ' + nrOfLicensesFound + ' licenses found.')
} else {
    console.log('Check complete. No licenses found.')
}