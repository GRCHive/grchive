<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleInfraUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class InfraTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            // Filter purposes
            name: `${inp.Name} ${inp.Purpose}`,
            value: inp
        }
    }

    goToInfra(item : any) {
        window.location.assign(createSingleInfraUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToInfra
                }
            }
        )
    }

}

</script>
