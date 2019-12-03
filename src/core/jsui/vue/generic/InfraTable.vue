<template>
    <base-resource-table
        :resources="resources"
        :value="value"
        :selectable="selectable"
        :multi="multi"
        :search="search"
        :table-headers="tableHeaders"
        :table-items="tableItems"
        @input="$emit('input', ...arguments)"
        @click:row="goToInfra"
    >
    </base-resource-table>
</template>

<script lang="ts">

import Vue from 'vue'
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
}

</script>
