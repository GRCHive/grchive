<script lang="ts">

import Vue, { VNode } from 'vue'
import { VBtn, VIcon, VDataTable } from 'vuetify/lib'
import Component, { mixins } from 'vue-class-component'
import ResourceTableProps from './ResourceTableProps'

const TableProps = Vue.extend({
    props: {
        tableHeaders: Array,
        tableItems: Array,
    }
})

@Component
export default class BaseResourceTable extends mixins(ResourceTableProps, TableProps) {
    selected: any[] = []

    get selectedSet() : Set<any> {
        return new Set<any>(this.selected)
    }

    changeInput(items: any[]) {
        this.$emit('input', items)
    }

    get finalTableHeaders() : any[] {
        let headers = this.tableHeaders

        if (this.useCrudDelete) {
            headers.push({
                text: 'Actions',
                value: 'action',
                sortable: false,
                filterable: false
            })
        }

        return headers
    }

    clickRow(item : any) {
        if (this.selectable) {
            if (this.selectedSet.has(item)) {
                this.selected = this.selected.filter((ele : any) => ele != item)
            } else {
                this.selected = [...this.selected, item]
            }

            this.changeInput(this.selected)
        } else {
            this.$emit('click:row', item)
        }
    }

    deleteItem(e : MouseEvent, item : any) {
        this.$emit('delete', item)
        e.stopPropagation()
    }

    renderActionSlot(props: any) : VNode {
        let childEle : VNode[] = []

        if (this.useCrudDelete) {
            childEle.push(
                this.$createElement(
                    VBtn,
                    {
                        props: {
                            icon: true,
                            'x-small': true,
                        },
                        on: {
                            click: (e : MouseEvent) => this.deleteItem(e, props.item)
                        }
                    },
                    [
                        this.$createElement(
                            VIcon,
                            'mdi-delete'
                        )
                    ]
                )
            )
        }

        return this.$createElement(
            'div',
            childEle
        )
    }

    render() : VNode {
        return this.$createElement(
            VDataTable,
            {
                props: {
                    value: this.selected,
                    headers: this.finalTableHeaders,
                    items: this.tableItems,
                    showSelect: this.selectable,
                    singleSelect: this.multi,
                    search: this.search
                },
                on: {
                    input: this.changeInput,
                    'click:row': this.clickRow
                },
                scopedSlots: {
                    ...this.$scopedSlots,
                    'item.action': this.renderActionSlot
                }
            }
        )
    }
}

</script>

<style scoped>

>>>tr {
    cursor: pointer !important;
}

</style>
