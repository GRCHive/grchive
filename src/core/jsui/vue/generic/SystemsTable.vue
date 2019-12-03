<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleSystemUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class SystemsTable extends ResourceTableProps {
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

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
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
        window.location.assign(createSingleSystemUrl(
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
                    'click:row': this.goToSystem
                },
                scopedSlots: {
                    'item.name': this.renderName,
                }
            }
        )
    }


}

</script>
