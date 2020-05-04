package grchive.core.data.sources

import grchive.core.data.types.grchive.ClientData

/**
 * A generic connection to some underlying raw data.
 */
open class RawDataSource(internal val clientData : ClientData) {
    val clientDataId get() = clientData.id
}
