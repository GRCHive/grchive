<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleDocCatUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class SystemsTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Category',
                value: 'category',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            category: `${inp.Name} ${inp.Description}`,
            name: inp.Name,
            value: inp
        }
    }

    goToDocumentation(item : any) {
        window.location.assign(createSingleDocCatUrl(PageParamsStore.state.organization!.OktaGroupName, item.value.Id))
    }

    renderCategory(props: any) : VNode { 
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
                    resourceName: "documentation categories"
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToDocumentation
                },
                scopedSlots: {
                    'item.category': this.renderCategory,
                }
            }
        )
    }
}

</script>

