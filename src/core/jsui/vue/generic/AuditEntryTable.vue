<script lang="ts">

import Vue, { VNode } from 'vue'
import Component, { mixins } from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { VTooltip, VIcon } from 'vuetify/lib'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import { createUserString } from '../../ts/users'
import { ResourceHandle, standardizeResourceType } from '../../ts/resourceUtils'
import { 
    TGetAuditTrailOutput, getAuditTrail,
    TAllAuditTrailOutput, allAuditTrail,
} from '../../ts/api/apiAuditTrail'
import { contactUsUrl } from '../../ts/url'
import { AuditEventEntry } from '../../ts/auditTrail'
import MetadataStore from '../../ts/metadata'

const Props = Vue.extend({
    props: {
        retrievalParams : {
            type: Object,
            default: () => Object(),
        }
    }
})

@Component({
    components: {
        BaseResourceTable
    }
})
export default class AuditEntryTable extends mixins(ResourceTableProps, Props) {
    eventIdToResourceHandle : Record<number, ResourceHandle | null> = Object()
    eventIdProcessed : Set<number> = new Set<number>()

    loadedData : AuditEventEntry[] = []
    totalEntries: number = 0
    currentPage: number = 1
    itemsPerPage: number = 10
    sortHeaders: string[] = []
    sortDesc: boolean[] = []

    isLoading : boolean = false

    @Watch('sortDesc')
    @Watch('sortHeaders')
    @Watch('itemsPerPage')
    @Watch('currentPage')
    @Watch('retrievalParams', { deep: true })
    refreshData() {
        this.isLoading = true

        let params = {
            orgId: PageParamsStore.state.organization!.Id,
            ...this.retrievalParams,
            page: this.currentPage - 1,
            numItems: this.itemsPerPage,
        }

        if (this.sortHeaders.length == 1) {
            params.sortHeaders = this.sortHeaders[0]
            if (this.sortDesc.length == 1 && this.sortDesc[0]) {
                params.sortDesc = true
            } else {
                params.sortDesc = false
            }
        }

        allAuditTrail(
            params
        ).then((resp : TAllAuditTrailOutput) => {
            this.loadedData = resp.data.Entries
            this.totalEntries = resp.data.Total
            this.isLoading = false
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

    get tableHeaders() : any[] {
        return [
            {
                text: 'Time',
                value: 'time',
            },
            {
                text: 'User',
                value: 'user',
            },
            {
                text: 'Action',
                value: 'gaction',
            },
            {
                text: 'Type',
                value: 'type',
            },
            {
                text: 'Resource',
                value: 'resource',
                sortable: false,
            },
        ]
    }

    get tableItems(): any[] {
        return this.loadedData.map(this.transformInputResourceToTableItem)
    }

    renderResource(props: any) : VNode { 
        this.retrieveAuditTrailDetails(props.item.value)

        if (props.item.id in this.eventIdToResourceHandle) {
            let res : ResourceHandle | null = this.eventIdToResourceHandle[props.item.id]
            if (!!res) {
                if (!!res.resourceUri) {
                    return this.$createElement(
                        'a',
                        {
                            attrs: {
                                href: res.resourceUri
                            }
                        },
                        res.displayText,
                    )
                } else {
                    return this.$createElement(
                        'span',
                        [
                            this.$createElement(
                                's',
                                res.displayText,
                            ),
                            this.$createElement(
                                VTooltip,
                                {
                                    props: {
                                        bottom: true,
                                        left: true,
                                    },
                                    scopedSlots: {
                                        activator: (props : any) : VNode => {
                                            return this.$createElement(
                                                VIcon,
                                                {
                                                    props: {
                                                        small: true
                                                    },
                                                    on: props.on
                                                },
                                                "mdi-help-circle"
                                            )
                                        }
                                    }
                                },
                                "This resource can not be found. It may have been deleted."
                            ),
                        ]
                    )
                }

            } else {
                return this.$createElement('span', 'Unknown')
            }
        } else {
            return this.$createElement('span', 'Loading...')
        }
    }

    retrieveAuditTrailDetails(inp : any) {
        if (!(inp.Id in this.eventIdToResourceHandle) && !this.eventIdProcessed.has(inp.Id)) {
            this.eventIdProcessed.add(inp.Id)
            getAuditTrail({
                orgId: PageParamsStore.state.organization!.Id,
                resourceHandleOnly: true,
                entryId: inp.Id,
            }).then((resp : TGetAuditTrailOutput) => {
                Vue.set(this.eventIdToResourceHandle, inp.Id, resp.data.Handle!)
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
                this.eventIdToResourceHandle[inp.Id] = null
            })
        }
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            gaction: inp.Action,
            time: standardFormatTime(inp.PerformedAt),
            type: standardizeResourceType(inp.ResourceType),
            user: createUserString(MetadataStore.getters.getUser(inp.UserId)),
            resource: this.eventIdToResourceHandle[inp.Id],
            value: inp
        }
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    resourceName: "audit trail entry",
                    showExpand: false,
                    loading: this.isLoading,
                    serverItemsLength: this.totalEntries,
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'update:page': (n : number) => { this.currentPage = n  },
                    'update:sort-by': (s : string[]) => { this.sortHeaders = s },
                    'update:sort-desc': (d : boolean[]) => { this.sortDesc = d },
                    'update:items-per-page': (n : number) => { this.itemsPerPage = n },
                },
                scopedSlots: {
                    'item.resource': this.renderResource,
                }
            }
        )
    }

    mounted() {
        this.refreshData()
    }
}

</script>
