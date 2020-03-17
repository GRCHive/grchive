export function standardFormatTime(dt : Date | null) : string {
    if (!dt) {
        return "None"
    }
    return dt.toLocaleString(undefined, {
        //@ts-ignore
        dateStyle: "short",
        timeZoneName: "short",
        timeStyle: "short"
    })
}

export function standardFormatDate(dt : Date | null) : string {
    if (!dt) {
        return "None"
    }
    return `${dt.getFullYear()}-${(dt.getMonth()+1).toString().padStart(2, "0")}-${dt.getDate().toString().padStart(2, "0")}`
}

export function createLocalDateFromDateString(str : string) : Date {
    let data = str.split('-')

    if (data.length != 3) {
        return new Date()
    }

    let dt = new Date()
    dt.setFullYear(
        parseInt(data[0]),
        parseInt(data[1])-1,
        parseInt(data[2])
    )
    return dt
}
