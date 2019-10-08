export function standardFormatTime(dt : Date) : string {
    return dt.toLocaleString(undefined, {
        //@ts-ignore
        dateStyle: "short",
        timeZoneName: "short",
        timeStyle: "short"
    })
}
