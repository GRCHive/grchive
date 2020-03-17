<script lang="ts">

import Vue, {VNode} from 'vue'
import { VIcon, VRow, VCol, VTooltip } from 'vuetify/lib'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { contactUsUrl, createSingleSqlRequestUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import { DbSqlQuery, DbSqlQueryMetadata, DbSqlQueryRequestApproval } from '../../ts/sql'
import { getSqlQuery, TGetSqlQueryOutput } from '../../ts/api/apiSqlQueries'
import {
    statusSqlRequest, TStatusSqlRequestOutput
} from '../../ts/api/apiSqlRequests'
import SqlTextArea from './SqlTextArea.vue'

interface QueryMetadataPacket {
    query: DbSqlQuery
    metadata : DbSqlQueryMetadata
}

@Component({
    components: {
        BaseResourceTable,
        SqlTextArea
    }
})
export default class SqlRequestTable extends ResourceTableProps {
    idToQuery : Record<number, QueryMetadataPacket> = Object()
    pendingQueryRequests : Set<number> = new Set<number>()

    idToApproval : Record<number, DbSqlQueryRequestApproval | null> = Object()
    pendingApprovalRequests : Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Query',
                value: 'query',
            },
            {
                text: 'Requester',
                value: 'requester',
            },
            {
                text: 'Request Time',
                value: 'requestTime',
                sort: (a : Date, b : Date) => {
                    return a.getTime() - b.getTime()
                }
            },
            {
                text: 'Assignee',
                value: 'assignee',
            },
            {
                text: 'Approval',
                value: 'approval',
            },
        ]
    }

    loadQuery(id : number) {
        if (id in this.idToQuery) {
            return
        }

        if (this.pendingQueryRequests.has(id)) {
            return
        }

        this.pendingQueryRequests.add(id)
        new Promise<QueryMetadataPacket>((resolve, reject) => {
            getSqlQuery({
                metadataId: -1,
                orgId: PageParamsStore.state.organization!.Id,
                queryId: id,
            }).then((resp : TGetSqlQueryOutput) => {
                resolve({
                    query: resp.data.Queries[0],
                    metadata: resp.data.Metadata,
                })
            }).catch((err : any) => {
                reject()
            })
        }).then((val : QueryMetadataPacket) => {
            Vue.set(this.idToQuery, id, val)
            this.pendingQueryRequests.delete(id)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    loadApproval(id : number) {
        if (id in this.idToApproval) {
            return
        }

        if (this.pendingApprovalRequests.has(id)) {
            return
        }

        this.pendingApprovalRequests.add(id)
        new Promise<DbSqlQueryRequestApproval | null>((resolve, reject) => {
            statusSqlRequest({
                requestId: id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp: TStatusSqlRequestOutput) => {
                resolve(resp.data)
            }).catch(() => {
                reject()
            })
        }).then((val : DbSqlQueryRequestApproval | null) => {
            Vue.set(this.idToApproval, id, val)
            this.pendingApprovalRequests.delete(id)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToSqlRequest(item : any) {
        window.location.assign(createSingleSqlRequestUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    transformInputResourceToTableItem(inp : any) : any {
        this.loadQuery(inp.QueryId)
        this.loadApproval(inp.Id)

        let queryName = !!this.idToQuery[inp.QueryId] ? 
            `${this.idToQuery[inp.QueryId].metadata.Name} v${this.idToQuery[inp.QueryId].query.Version}` :
            "Loading..."

        return {
            id: inp.Id,
            name: inp.Name,
            query: queryName,
            requester: createUserString(MetadataStore.getters.getUser(inp.UploadUserId)),
            requestTime: inp.UploadTime,
            approval: this.idToApproval[inp.Id],
            assignee: createUserString(MetadataStore.getters.getUser(inp.AssigneeUserId)),
            value: inp
        }
    }

    renderRequestTime(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.requestTime)
        )
    }

    renderApproval(props : any) : VNode {
        let approval : DbSqlQueryRequestApproval | null = props.item.approval

        if (!approval) {
            return this.$createElement(
                VIcon,
                {
                    props: {
                        small : true,
                        color: 'warning'
                    }
                },
                'mdi-help-circle'
            )
        } else if (approval.Response) {
            return this.$createElement(
                VIcon,
                {
                    props: {
                        small : true,
                        color: 'success'
                    }
                },
                'mdi-check'
            )
        } else {
            let renderIcon = (props : any) => {
                return this.$createElement(
                    VIcon,
                    {
                        props: {
                            small : true,
                            color: 'error'
                        },
                        on: props.on,
                    },
                    'mdi-close'
                )
            }

            return this.$createElement(
                VTooltip,
                {
                    props: {
                        bottom: true,
                    },
                    scopedSlots: {
                        activator: renderIcon
                    }
                },
                approval.Reason
            )
        }
    }

    renderExpansion(props : any) : VNode {
        let descriptionCol = this.$createElement(
            VCol,
            {
                props: {
                    cols: 4,
                }
            },
            props.item.value.Description
        )

        let queryColChildren : VNode[] = []

        if (!!this.idToQuery[props.item.value.QueryId])  {
            queryColChildren.push(
                this.$createElement(
                    SqlTextArea,
                    {
                        props: {
                            value: this.idToQuery[props.item.value.QueryId].query.Query,
                            readonly: true,
                        },
                    },
                )
            )
        }

        let queryCol = this.$createElement(
            VCol,
            {
                props: {
                    cols: 8,
                }
            },
            queryColChildren
        )

        return this.$createElement(
            'td',
            {
                attrs: {
                    colspan: props.headers.length
                },
            },
            [this.$createElement(
                VRow,
                {
                    props: {
                        align: "center",
                        justify: "center",
                    }
                },
                [
                    descriptionCol, queryCol,
                ]
            )]
        )
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToSqlRequest
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.approval': this.renderApproval,
                    'item.requestTime': this.renderRequestTime,
                }
            }
        )
    }
}

</script>
