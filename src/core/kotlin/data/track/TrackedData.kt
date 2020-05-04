package grchive.core.data.track

class TrackedData<T> {
    val t : T
    internal val source : TrackedSource

    constructor(inT : T, inSource : TrackedSource) {
        t = inT
        source = inSource
        source.addData(this)
    }

    override fun equals(other : Any?) : Boolean {
        if (other == null) {
            return false
        }
        return (other is TrackedData<*> && t == other.t)
    }

    override fun hashCode() : Int {
        return t.hashCode()
    }
}
