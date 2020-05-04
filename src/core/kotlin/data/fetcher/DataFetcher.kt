package grchive.core.data.fetcher

import grchive.core.data.sources.RawDataSource
import grchive.core.data.filters.Filter
import grchive.core.data.track.TrackedData

interface DataFetcher<T, D> {
    /**
     * A function to obtain all the data that maches the given filter(s).
     *
     * @param source The [RawDataSource] to retrieve data from.
     * @param filters A map of filters to use to filter data.
     * @return All the matching data of type T.
     */
    fun fetch(source : D, filters : Map<String, Filter>) : List<TrackedData<T>>
}
