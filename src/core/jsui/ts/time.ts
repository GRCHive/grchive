export enum Days {
    Sunday = 0,
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
}

export type DaysKey = keyof typeof Days

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

export function standardFormatTimeOnly(dt : Date | null ) : string {
    if (!dt) {
        return "None"
    }

    const amPm = (dt.getHours() >= 12) ? 'PM' : 'AM'
    const hours = (dt.getHours() == 0 || dt.getHours() == 12) ? '12' : `${dt.getHours() % 12}`
    return `${hours}:${dt.getMinutes().toString().padStart(2, "0")} ${amPm} UTC${dt.getTimezoneOffset() > 0 ? '-' : '+'}${Math.abs(dt.getTimezoneOffset()) / 60}`
}

export function standardFormatDate(dt : Date | null) : string {
    if (!dt) {
        return "None"
    }
    return `${dt.getFullYear()}-${(dt.getMonth()+1).toString().padStart(2, "0")}-${dt.getDate().toString().padStart(2, "0")}`
}

export function vuetifyCalendarTimeFormat(dt : Date) : string {
    return `${standardFormatDate(dt)} ${dt.getHours()}:${dt.getMinutes()}`
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

export let DaysSelectItems = Object.keys(Days)
    .filter((key : any) => !isNaN(Number(Days[key])))
    .map((key : any) => ({
        text: key,
        value: Days[key]
    }))

export interface TimeRange {
    Start : Date
    End: Date
}

export function cleanTimeRangeFromJson(r : TimeRange) {
    r.Start = new Date(r.Start)
    r.End = new Date(r.End)
}
