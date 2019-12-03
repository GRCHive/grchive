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
        @click:row="goToSystem"
    >
        <template v-slot:item.name="{ item }">
            <p class="ma-0 pa-0 body-1 font-weight-bold">{{ item.value.Name }}</p>
            <p class="ma-0 pa-0 caption font-weight-light">{{ item.value.Description }}</p>
        </template>
    </base-resource-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleSystemUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class SystemsTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Purpose',
                value: 'purpose',
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
            name: `${inp.Name} ${inp.Description}`,
            purpose: inp.Purpose,
            value: inp
        }
    }

    transformTableItemToInputResource(inp : any) : any {
        return inp.value
    }

    goToSystem(item : any) {
        window.location.assign(createSingleSystemUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }
}

</script>
