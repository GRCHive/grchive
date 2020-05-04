package grchive.core.test.comparisons

import grchive.core.data.track.TrackedData
import grchive.core.test.TestAction
import grchive.core.test.TestContainer

import kotlin.reflect.KProperty1
import kotlin.reflect.jvm.javaField

fun <T> testEqual(cn : TestContainer, a : TrackedData<T>?, b : TrackedData<T>?) : Boolean {
    val ok = (a == b)
    cn.logTest(TestAction.Equal, a, b, ok)
    return ok
}

fun <T1,T2> testFieldEqual(cn : TestContainer, a : TrackedData<T1>?, b : TrackedData<T2>?, propa : KProperty1<T1, *>, propb : KProperty1<T2, *>) : Boolean {
    val at = a?.t
    val bt = b?.t

    var ok : Boolean
    if ((at == null && bt != null) || (at != null && bt == null) || (at == null && bt == null)) {
        ok = false
    } else {
        ok = propa.get(at!!) == propb.get(bt!!)
    }

    cn.logTest(TestAction.Equal, a, b, ok, "${propa.javaField?.getName()} -> ${propb.javaField?.getName()}")
    return ok
}
