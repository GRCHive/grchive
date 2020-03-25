package grchive.core.data.types.grchive

import java.time.OffsetDateTime

/**
 * Contains basic information about a process flow.
 *
 * @property id The underlying unique identifier for the process flow.
 * @property name A human readable name.
 * @property orgId Id of the organization this process flow belongs to.
 * @property description An optional description.
 * @property creationTime When the process flow was first created.
 * @property lastUpdatedTime When the process flow metadata was last updated.
 */
data class ProcessFlowMetadata (
    val id : Long,
    val name : String,
    val orgId: Int,
    val description : String,
    val creationTime : OffsetDateTime,
    val lastUpdatedTime : OffsetDateTime
)
