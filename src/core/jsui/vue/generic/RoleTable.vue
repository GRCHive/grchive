<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import { VCheckbox } from 'vuetify/lib'
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

    renderRoleName(props : any): VNode {
        //<template v-slot:item.roleName="{ item }">
        //    <p class="ma-0 pa-0 body-1 font-weight-bold">{{ item.value.Name }}</p>
        //    <p class="ma-0 pa-0 caption font-weight-light">{{ item.value.Description }}</p>
        //</template>
        return this.$createElement(
            'div',
            [
                this.$createElement(
                    'p',
                    {
                        class: {
                            'ma-0': true,
                            'pa-0': true,
                            'body-1': true,
                            'font-weight-bold': true
                        },
                        domProps: {
                            innerHTML: props.item.value.Name
                        }
                    },
                ),
                this.$createElement(
                    'p',
                    {
                        class: {
                            'ma-0': true,
                            'pa-0': true,
                            'caption': true,
                            'font-weight-light': true
                        },
                        domProps: {
                            innerHTML: props.item.value.Description
                        }
                    }
                ),
            ]
        )
    }

    renderIsDefault(props : any): VNode {
        //<template v-slot:item.isDefault="{ item }">
        //    <v-checkbox class="ma-0 pa-0" :input-value="item.value.IsDefault" disabled hide-details></v-checkbox>
        //</template>
        return this.$createElement(
            VCheckbox,
            {
                class: {
                    'ma-0': true,
                    'pa-0': true,
                },
                props: {
                    'input-value': props.item.value.IsDefault,
                    disabled: true,
                    'hide-details': true
                }
            }
        )
    }

    renderIsAdmin(props : any): VNode {
        //<template v-slot:item.isAdmin="{ item }">
        //    <v-checkbox class="ma-0 pa-0" :input-value="item.value.IsAdmin" disabled hide-details></v-checkbox>
        //</template>
        return this.$createElement(
            VCheckbox,
            {
                class: {
                    'ma-0': true,
                    'pa-0': true,
                },
                props: {
                    'input-value': props.item.value.IsAdmin,
                    disabled: true,
                    'hide-details': true
                }
            }
        )
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
                    input: (...args : any[]) => this.$emit('input', ...args),
                    'click:row': this.goToRole
                },
                scopedSlots: {
                    'item.roleName': this.renderRoleName,
                    'item.isDefault': this.renderIsDefault,
                    'item.isAdmin': this.renderIsAdmin
                }
            }
        )
    }

}

</script>
