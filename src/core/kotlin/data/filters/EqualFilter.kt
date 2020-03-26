package grchive.core.data.filters

final class EqualFilter<T>(val target : T) : Filter {
    override fun createSqlCondition(col : String) : Pair<String, ArrayList<*>> {
        return "${col} = ?" to arrayListOf(target)
    }
}
