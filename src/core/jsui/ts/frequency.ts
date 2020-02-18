export const frequencyTypes : string[] = [
    "Days",
    "Weeks",
    "Months",
    "Quarters",
    "Years"
]

export function createFrequencyDisplayString(freqType: number, freqInterval: number) : string {
    if (freqType == -1) {
        return "Ad-Hoc"
    }
    return `Every ${freqInterval} ${frequencyTypes[freqType]}`
}
