<template>
    <v-data-table
        v-if="ready"
        v-model="selected"
        :headers="tableHeaders"
        :items="tableItems"
        :show-select="selectable"
        :single-select="!multi"
        :search="search"
        @input="changeInput"
        @click:row="goToDb">
    </v-data-table>
</template>

<script lang="ts">

import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleDbUrl } from '../../ts/url'
import { getDbTypeAsString } from '../../ts/databases'
import MetadataStore from '../../ts/metadata'

@Component
export default class SystemsTable extends BaseResourceTable {
    get ready() : boolean {
        return MetadataStore.state.dbTypesInitialized
    }

    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Type',
                value: 'type',
            },
            {
                text: 'Version',
                value: 'version',
            },
        ]
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: inp.Name,
            type: getDbTypeAsString(inp),
            version: inp.Version,
            value: inp
        }
    }

    transformTableItemToInputResource(inp : any) : any {
        return inp.value
    }

    goToDb(item : any) {
        window.location.assign(createSingleDbUrl(
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
