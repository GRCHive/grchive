export function isObjectEmpty(obj : any) : boolean {
    return Object.getOwnPropertyNames(obj).length === 0
}
