<template>
    <base-resource-table
        v-if="ready"
        :resources="resources"
        :value="value"
        :selectable="selectable"
        :multi="multi"
        :search="search"
        :table-headers="tableHeaders"
        :table-items="tableItems"
        @input="$emit('input', ...arguments)"
        @click:row="goToDb"
    >
    </base-resource-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleDbUrl } from '../../ts/url'
import { getDbTypeAsString } from '../../ts/databases'
import MetadataStore from '../../ts/metadata'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class SystemsTable extends ResourceTableProps {
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

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
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
