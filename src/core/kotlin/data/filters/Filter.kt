package grchive.core.data.filters

interface Filter {
    /**
     * Creates a conditional that can be used in a SQL statement.
     *
     * @param col The column to use in the check.
     * @return A Pair<String,ArrayList<T>> where the String is the conditional and 
     *         the ArrayList<T> contains any arguments that need to be bound to the SQL
     *         statement for this filter.
     */
    fun createSqlCondition(col : String) : Pair<String, ArrayList<*>>
}
