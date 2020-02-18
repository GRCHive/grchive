export const frequencyTypes : string[] = [
    "Days",
    "Weeks",
    "Months",
    "Quarters",
    "Years"
]

export function createFrequencyDisplayString(freqType: number, freqInterval: number, freqOther : string) : string {
    if (freqType == -1) {
        return "Ad-Hoc"
    } else if (freqType == -2) {
        return `Other: ${freqOther}`
    }
    return `Every ${freqInterval} ${frequencyTypes[freqType]}`
}
