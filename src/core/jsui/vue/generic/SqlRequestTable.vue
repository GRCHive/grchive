<script lang="ts">

import Vue, {VNode} from 'vue'
import { VIcon, VRow, VCol } from 'vuetify/lib'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import { DbSqlQuery, DbSqlQueryMetadata } from '../../ts/sql'
import { getSqlQuery, TGetSqlQueryOutput } from '../../ts/api/apiSqlQueries'
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
    pendingRequests : Set<number> = new Set<number>()

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
            },
        ]
    }

    loadQuery(id : number) {
        if (id in this.idToQuery) {
            return
        }

        if (this.pendingRequests.has(id)) {
            return
        }

        this.pendingRequests.add(id)
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
            this.pendingRequests.delete(id)
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
        window.location.assign('#')
    }

    transformInputResourceToTableItem(inp : any) : any {
        this.loadQuery(inp.QueryId)

        let queryName = !!this.idToQuery[inp.QueryId] ? 
            `${this.idToQuery[inp.QueryId].metadata.Name} v${this.idToQuery[inp.QueryId].query.Version}` :
            "Loading..."

        return {
            id: inp.Id,
            name: inp.Name,
            query: queryName,
            requester: createUserString(MetadataStore.getters.getUser(inp.UploadUserId)),
            requestTime: standardFormatTime(inp.UploadTime),
            value: inp
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
                }
            }
        )
    }
}

</script>
