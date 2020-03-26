package grchive.core.data.fetcher

interface DataFetcher<T> {
    fun fetch() : T
}
