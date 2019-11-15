<template>
    <v-data-table
        v-model="selected"
        :headers="tableHeaders"
        :items="tableItems"
        :show-select="selectable"
        :single-select="!multi"
        :search="search"
        @input="changeInput">

        <template v-slot:item.roleName="{ item }">
            <p class="ma-0 pa-0 body-1 font-weight-bold">{{ item.value.Name }}</p>
            <p class="ma-0 pa-0 caption font-weight-light">{{ item.value.Description }}</p>
        </template>

        <template v-slot:item.isDefault="{ item }">
            <v-checkbox class="ma-0 pa-0" :input-value="item.value.IsDefault" disabled hide-details></v-checkbox>
        </template>

        <template v-slot:item.isAdmin="{ item }">
            <v-checkbox class="ma-0 pa-0" :input-value="item.value.IsAdmin" disabled hide-details></v-checkbox>
        </template>

        <template v-slot:item.action="{ item }">
            <v-btn small icon :href="createOrgRoleUrl(orgGroupName, item.value.Id)">
                <v-icon small>mdi-pencil</v-icon>
            </v-btn>
        </template>
    </v-data-table>
</template>

<script lang="ts">

import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import { RoleMetadata } from '../../ts/roles'
import { PageParamsStore } from '../../ts/pageParams'
import { createOrgRoleUrl } from '../../ts/url'

@Component({
    methods: {
        createOrgRoleUrl
    }
})
export default class RoleTable extends BaseResourceTable {
    get orgGroupName() : string {
        return PageParamsStore.state.organization!.OktaGroupName
    }

    get tableHeaders() : any[] {
        return [
            {
                text: 'Role',
                value: 'roleName',
            },
            {
                text: 'Default',
                value: 'isDefault',
                filterable: false,
            },
            {
                text: 'Admin',
                value: 'isAdmin',
                filterable: false,
            },
            {
                text: 'Actions',
                value: 'action',
                sortable: false,
                filterable: false,
            },
        ]
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            // Need this for the filter
            roleName: `${inp.Name} ${inp.Description}`,
            value: inp
        }
    }

    transformTableItemToInputResource(inp : any) : any {
        return inp.value
    }
}

</script>
