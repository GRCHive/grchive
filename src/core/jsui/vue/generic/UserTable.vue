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
    >
    </base-resource-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class UserTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'fullName'
            },
            {
                text: 'Email',
                value: 'email'
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            fullName: `${inp.FirstName} ${inp.LastName}`,
            email: inp.Email,
            value: inp
        }
    }

    transformTableItemToInputResource(inp : any) : any {
        return inp.value
    }
}

</script>
