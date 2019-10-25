// This map will store relationship data between TWO "things".
// Let us call these things Thing A and Thing B.
// Then this map will allow us to know:
//  1) All thing B's that are related to a thing A.
//  2) All thing A's that are related to a thing B.
//  3) Whether or not a thing A and thing B are related.
// This map will also faciliate adding/editing/deleting these relationships as well.
class RelationshipMap<A=any, B=any, T=void> {
    aToB : Map<A, Set<B>>
    bToA : Map<B, Set<A>>
    dataStore : Map<A, Map<B, T>>
    changed : 1

    constructor() {
        this.aToB = new Map<A, Set<B>>()
        this.bToA = new Map<B, Set<A>>()
        this.dataStore = new Map<A, Map<B,T>>()
        this.changed = 1
    }

    add(a : A, b : B, v : T) {
        if (!this.aToB.has(a)) {
            this.aToB.set(a, new Set<B>()) 
        }
        this.aToB.get(a)!.add(b)

        if (!this.bToA.has(b)) {
            this.bToA.set(b, new Set<A>()) 
        }
        this.bToA.get(b)!.add(a)

        if (!this.dataStore.has(a)) {
            this.dataStore.set(a, new Map<B,T>())
        }

        let substore : Map<B,T> = this.dataStore.get(a)!
        substore.set(b, v)
        this.changed += 1
    }

    delete(a : A, b: B) {
        if (this.aToB.has(a)) {
            this.aToB.get(a)!.delete(b)
        }

        if (this.bToA.has(b)) {
            this.bToA.get(b)!.delete(a)
        }

        if (this.dataStore.has(a)) {
            let substore = this.dataStore.get(a)!
            if (substore.has(b)) {
                substore.delete(b)
            }
        }
        this.changed += 1
    }

    deleteA(a : A) {
        this.aToB.delete(a)
        this.bToA.forEach((value, _) => {
            value.delete(a)
        })

        if (this.dataStore.has(a)) {
            this.dataStore.delete(a)
        }
        this.changed += 1
    }

    deleteB(b : B) {
        this.aToB.forEach((value, _) => {
            value.delete(b)
        })
        this.bToA.delete(b)

        this.dataStore.forEach((value, _) => {
            value.delete(b)
        })
        this.changed += 1
    }

    getB(a : A) : B[] {
        if (!this.aToB.has(a)) {
            return []
        }

        return Array.from(this.aToB.get(a)!)
    }

    getA(b : B) : A[] {
        if (!this.bToA.has(b)) {
            return []
        }
        return Array.from(this.bToA.get(b)!)
    }

    touch() {
        this.changed += 1
    }
}

export default RelationshipMap
