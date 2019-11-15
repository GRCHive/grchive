export function setSubtract<T>(a : Set<T>, b : Set<T>) : Set<T> {
    return new Set<T>([...a].filter((x : T) => !b.has(x)))
}

export function listSubtract<T>(a : T[], b : T[]) : T[] {
    let sa = new Set<T>(a)
    let sb = new Set<T>(b)
    let sc = setSubtract(sa, sb)
    return [...sc]
}
