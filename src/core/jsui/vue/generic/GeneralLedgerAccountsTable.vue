<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleGLAccountUrl } from '../../ts/url'
import { getGLAccountParentString } from '../../ts/generalLedger'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class GLAccountsTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Parent',
                value: 'parent',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            // Filter purposes
            name: `${inp.AccountName} ${inp.AccountId} ${inp.AccountDescription}`,
            parent: getGLAccountParentString(inp),
            value: inp
        }
    }

    goToGLAccount(item : any) {
        window.location.assign(createSingleGLAccountUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }

    renderName(props : any): VNode {
        //<template v-slot:item.name="{ item }">
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
                            innerHTML: `${props.item.value.AccountName} (${props.item.value.AccountId})`
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
                            innerHTML: props.item.value.AccountDescription
                        }
                    }
                ),
            ]
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
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToGLAccount
                },
                scopedSlots: {
                    'item.name': this.renderName,
                }
            }
        )
    }
}

</script>
