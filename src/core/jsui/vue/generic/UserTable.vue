<script lang="ts">

import Vue, { VNode } from 'vue'
import Component, { mixins } from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { RoleMetadata } from '../../ts/roles'

const RoleProps = Vue.extend({
    props: {
        availableRoles: {
            type: Object as () => Record<number, RoleMetadata>,
            default: Object()
        },
        showRole: {
            type: Boolean,
            default: false
        }
    }
})

@Component({
    components: {
        BaseResourceTable
    }
})
export default class UserTable extends mixins(ResourceTableProps, RoleProps) {
    get tableHeaders() : any[] {
        let headers = [
            {
                text: 'Name',
                value: 'fullName'
            },
            {
                text: 'Email',
                value: 'email'
            },
        ]

        if (this.showRole) {
            headers.push({
                text: 'Role',
                value: 'role'
            })
        }

        return headers
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            fullName: `${inp.FirstName} ${inp.LastName}`,
            email: inp.Email,
            value: inp,
            role: this.showRole ? 
                this.availableRoles[inp.RoleId].Name :
                undefined
        }
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                }
            }
        )
    }
}

</script>
