package grchive.core.test.matchers

import grchive.core.data.track.TrackedData
import grchive.core.test.TestContainer

import kotlin.reflect.KProperty1

fun <T1,T2,K> findMatchesUsingKey(
    a : List<TrackedData<T1>>,
    b : List<TrackedData<T2>>,
    keya : KProperty1<T1, K>,
    keyb : KProperty1<T2, K>,
    compare : (a : TrackedData<T1>?, b : TrackedData<T2>?) -> Unit
) {
    val aMap = a.map { keya.get(it.t) to it }.toMap()
    val bMap = b.map { keyb.get(it.t) to it }.toMap()
    val processed = mutableSetOf<K>()

    aMap.forEach loop@{
        k, v -> 
            if (processed.contains(k)) {
                return@loop
            }
            processed.add(k)

            val other = bMap.get(k)
            compare(v, other)
    }

    bMap.forEach loop@{
        k, v -> 
            if (processed.contains(k)) {
                return@loop
            }
            processed.add(k)

            val other = aMap.get(k)
            compare(other, v)
    }

}
