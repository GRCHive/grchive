export class TRange<T> {
    min : T
    max : T

    constructor(min : T, max : T) {
        this.min = min
        this.max = max
    }

    intersects(other : TRange<T>) : boolean {
        if (this.max < other.min || this.min > other.max) {
            return false
        }
        return true
    }
}
