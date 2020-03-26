package grchive.core.internal.database

import java.lang.StringBuilder

import grchive.core.data.filters.Filter

fun createSqlConditionsFromFilters(filters : Map<String, Filter>, prefix: String = "", columnPrefix : String = "") : Pair<String,ArrayList<*>> {
    if (filters.isEmpty()) {
        return "" to arrayListOf<Any?>()
    }

    val sb = StringBuilder()
    if (prefix.isNotBlank()) {
        sb.append("${prefix} ")
    }

    val retArgs = arrayListOf<Any?>()
    var count = 0
    filters.forEach{
        k, v ->
            val (condition, args) = v.createSqlCondition("${columnPrefix}${k}")
            retArgs.addAll(args)

            sb.append(condition)
            if (count != filters.size - 1) {
                sb.append(" AND")
            }
            sb.append("\n")
            count += 1
    }
    return sb.toString() to retArgs
}
