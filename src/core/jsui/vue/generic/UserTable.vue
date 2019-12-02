<template>
    <v-data-table
        :value="selected"
        :headers="tableHeaders"
        :items="tableItems"
        :show-select="selectable"
        :single-select="!multi"
        :search="search"
        @input="changeInput">
    </v-data-table>
</template>

<script lang="ts">

import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import Vue from 'vue'

@Component
export default class UserTable extends BaseResourceTable {
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
