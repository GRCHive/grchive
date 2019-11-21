<template>
    <v-data-table
        v-model="selected"
        :headers="tableHeaders"
        :items="tableItems"
        :show-select="selectable"
        :single-select="!multi"
        :search="search"
        @input="changeInput"
        @click:row="goToInfra">
    </v-data-table>
</template>

<script lang="ts">

import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleInfraUrl } from '../../ts/url'

@Component
export default class InfraTable extends BaseResourceTable {
    get tableHeaders() : any[] {
        return [
        ]
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

<style scoped>

>>>tr {
    cursor: pointer !important;
}

</style>
