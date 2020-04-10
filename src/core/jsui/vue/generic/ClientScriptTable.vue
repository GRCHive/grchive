<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { ResourceHandle } from '../../ts/resourceUtils'
import { contactUsUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class ClientScriptTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: inp.Name,
            value: inp
        }
    }

    goToScript(item : any) {
        //window.location.assign(createSingleClientScriptUrl(
        //    PageParamsStore.state.organization!.OktaGroupName,
        //    item.value.Data.Id
        //))
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
                    resourceName: "scripts",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToScript
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                }
            }
        )
    }
}

</script>
