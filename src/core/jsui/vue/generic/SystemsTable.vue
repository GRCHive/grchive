<template>
    <v-data-table
        :value="selected"
        :headers="tableHeaders"
        :items="tableItems"
        :show-select="selectable"
        :single-select="!multi"
        :search="search"
        @input="changeInput"
        @click:row="goToSystem">

        <template v-slot:item.name="{ item }">
            <p class="ma-0 pa-0 body-1 font-weight-bold">{{ item.value.Name }}</p>
            <p class="ma-0 pa-0 caption font-weight-light">{{ item.value.Description }}</p>
        </template>
    </v-data-table>
</template>

<script lang="ts">

import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleSystemUrl } from '../../ts/url'

@Component
export default class SystemsTable extends BaseResourceTable {
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
        if (this.selectable) {
            this.manualToggleItem(item)
        } else {
            window.location.assign(createSingleSystemUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                item.value.Id))
        }
    }
}

</script>

<style scoped>

>>>tr {
    cursor: pointer !important;
}

</style>
