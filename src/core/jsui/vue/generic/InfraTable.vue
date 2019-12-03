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

    transformTableItemToInputResource(inp : any) : any {
        return inp.value
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
                    input: (...args : any[]) => this.$emit('input', ...args),
                    'click:row': this.goToInfra
                }
            }
        )
    }

}

</script>
