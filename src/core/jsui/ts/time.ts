export function standardFormatTime(dt : Date) : string {
    return dt.toLocaleString(undefined, {
        //@ts-ignore
        dateStyle: "short",
        timeZoneName: "short",
        timeStyle: "short"
    })
}

export function standardFormatDate(dt : Date) : string {
    return `${dt.getFullYear()}-${(dt.getMonth()+1).toString().padStart(2, "0")}-${dt.getDate().toString().padStart(2, "0")}`
}

export function createLocalDateFromDateString(str : string) : Date {
    let dt = new Date(str)
    dt.setDate(dt.getUTCDate())
    return dt
}
