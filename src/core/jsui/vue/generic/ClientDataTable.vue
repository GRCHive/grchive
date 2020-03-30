<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { ResourceHandle } from '../../ts/resourceUtils'
import { DataSourceLink } from '../../ts/clientData'
import { getDataSource, TGetDataSourceOutput } from '../../ts/api/apiDataSource'
import { contactUsUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class ClientDataTable extends ResourceTableProps {
    dataIdToSource : Record<number, ResourceHandle | null> = Object()
    dataIdProcessed : Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Source',
                value: 'source',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Data.Id,
            name: inp.Data.Name,
            source: this.dataIdToSource[inp.Data.Id],
            value: inp
        }
    }

    goToData(item : any) {
    }

    retrieveSourceResourceHandle(dataId : number, link : DataSourceLink) {
        if (!(dataId in this.dataIdToSource) && !this.dataIdProcessed.has(dataId)) {
            this.dataIdProcessed.add(dataId)
            getDataSource({
                source: link,
            }).then((resp : TGetDataSourceOutput) => {
                Vue.set(this.dataIdToSource, dataId, resp.data)
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
    }

    renderSource(props : any) : VNode {
        this.retrieveSourceResourceHandle(props.item.value.Data.Id, props.item.value.Link)

        let dataId = props.item.value.Data.Id
        if (dataId in this.dataIdToSource) {
            let handle = this.dataIdToSource[dataId]

            if (!!handle) {
                if (!!handle.resourceUri) {
                    return this.$createElement(
                        'a',
                        {
                            attrs: {
                                href: handle.resourceUri
                            }
                        },
                        handle.displayText,
                    )
                } else {
                    return this.$createElement('span', handle.displayText)
                }
            } else {
                return this.$createElement('span', 'Unknown')
            }
        } else {
            return this.$createElement('span', 'Loading...')
        }
    }

    renderExpansion(props : any) : VNode {
        return this.$createElement(
            'td',
            {
                attrs: {
                    colspan: props.headers.length
                },
            },
            props.item.value.Data.Description
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
                    resourceName: "data objects",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToData
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.source': this.renderSource,
                }
            }
        )
    }
}


</script>
