<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleVendorUrl } from '../../ts/url'
import { sanitizeUrl } from '../../ts/urlUtility'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class VendorTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Url',
                value: 'url',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            // Filter purposes
            name: inp.Name,
            url: inp.Url,
            value: inp
        }
    }

    goToVendor(item : any) {
        window.location.assign(createSingleVendorUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }

    renderUrl(props : any) : VNode {
        return this.$createElement(
            'a',
            {
                attrs: {
                    href: sanitizeUrl(props.item.url),
                    target: "_blank"
                },
                on: {
                    'click': (e : MouseEvent) => { e.stopPropagation() }
                }
            },
            props.item.url
        )
    }

    renderExpansion(props : any) : VNode {
        return this.$createElement(
            'td',
            {
                attrs: {
                    colspan: props.headers.length
                },
            },
            props.item.value.Description
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
                    resourceName: "vendors",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToVendor
                },
                scopedSlots: {
                    'item.url': this.renderUrl,
                    'expanded-item': this.renderExpansion
                }
            }
        )
    }
}

</script>
