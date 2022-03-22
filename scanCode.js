const scanResults = require('./sample.json')
let nrOfLicensesFound = 0

console.log('Checking ScanCode result for licenses')

for (var i = 0; i < scanResults.files.length; i++) {
    if (scanResults.files[i].licenses.length > 0) {
        console.log('License(s) found in file: ' + scanResults.files[i].base_name)
        for (var j = 0; j < scanResults.files[i].licenses.length; i++) {
            nrOfLicensesFound++
            console.log(scanResults.files[i].licenses[j])
        }
    }
}

if (nrOfLicensesFound > 0) {
    console.log('Check complete. ' + nrOfLicensesFound + 'licenses found.')
} else {
    console.log('Check complete. No licenses found.')
}