<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import { createUserString } from '../../ts/users'
import { getResourceHandle, ResourceHandle, standardizeResourceType } from '../../ts/resourceUtils'
import { contactUsUrl } from '../../ts/url'
import MetadataStore from '../../ts/metadata'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class AuditEntryTable extends ResourceTableProps {
    eventIdToResourceHandle : Record<number, ResourceHandle | null> = Object()

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
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    renderResource(props: any) : VNode { 
        console.log("RENDER RESOURCE:" , props.item, props.item.Id)
        if (props.item.id in this.eventIdToResourceHandle) {
            let res : ResourceHandle | null = this.eventIdToResourceHandle[props.item.id]
            if (!!res) {
                return this.$createElement(
                    'a',
                    {
                        attrs: {
                            href: res.resourceUri ? res.resourceUri : '#'
                        }
                    },
                    res.displayText,
                )

            } else {
                return this.$createElement('span', 'Unknown')
            }
        } else {
            return this.$createElement('span', 'Loading...')
        }
    }

    transformInputResourceToTableItem(inp : any) : any {
        if (!(inp.Id in this.eventIdToResourceHandle)) {
            getResourceHandle(inp.Action, inp.ResourceType, inp.ResourceId, inp.ResourceExtraData).then((resp : ResourceHandle | null) => {
                console.log("RESOLVE: ", inp.Id)
                Vue.set(this.eventIdToResourceHandle, inp.Id, resp)
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
                    showExpand: false
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                },
                scopedSlots: {
                    'item.resource': this.renderResource,
                }
            }
        )
    }
}

</script>
