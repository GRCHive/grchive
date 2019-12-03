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
        @click:row="goToRole"
    >
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
    </base-resource-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { RoleMetadata } from '../../ts/roles'
import { PageParamsStore } from '../../ts/pageParams'
import { createOrgRoleUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class RoleTable extends ResourceTableProps {
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
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
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

    goToRole(item : any) {
        window.location.assign(createOrgRoleUrl(PageParamsStore.state.organization!.OktaGroupName, item.value.Id))
    }
}

</script>
