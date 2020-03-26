package grchive.core.data.types.grchive

import java.time.OffsetDateTime
import org.jdbi.v3.core.mapper.reflect.ColumnName

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
    @ColumnName("id") val id : Long,
    @ColumnName("name") val name : String,
    @ColumnName("org_id") val orgId: Int,
    @ColumnName("description") val description : String,
    @ColumnName("created_time") val creationTime : OffsetDateTime,
    @ColumnName("last_updated_time") val lastUpdatedTime : OffsetDateTime
)
