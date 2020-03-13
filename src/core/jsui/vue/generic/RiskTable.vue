<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { contactUsUrl, createRiskUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'

@Component({
    components: {
        BaseResourceTable,
    }
})
export default class RiskTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToRisk(item : any) {
        window.location.assign(createRiskUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: `${inp.Name} ${inp.Identifier} ${inp.Description}`,
            value: inp
        }
    }

    renderName(props: any) : VNode { 
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
                            innerHTML: props.item.value.Identifier
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
                            innerHTML: props.item.value.Name
                        }
                    }
                ),
            ]
        )
    }

    renderExpansion(props : any) : VNode {
        return this.$createElement(
            'td',
            {
                attrs: {
                    colspan: props.headers.length
                },
            },
            props.item.value.Description
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
                    showExpand: true,
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any, globalDelete: boolean) => this.$emit('delete', item.value, globalDelete),
                    'click:row': this.goToRisk
                },
                scopedSlots: {
                    'item.name': this.renderName,
                    'expanded-item': this.renderExpansion,
                }
            }
        )
    }
}

</script>
