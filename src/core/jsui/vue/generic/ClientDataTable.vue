<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class ClientDataTable extends ResourceTableProps {
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
            source: '',
            value: inp
        }
    }

    goToData(item : any) {
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
                    'expanded-item': this.renderExpansion
                }
            }
        )
    }
}


</script>
