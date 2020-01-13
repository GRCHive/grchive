<script lang="ts">

import Vue, { VNode } from 'vue'
import { VBtn, VIcon, VDataTable, VDialog } from 'vuetify/lib'
import Component, { mixins } from 'vue-class-component'
import ResourceTableProps from './ResourceTableProps'
import GenericDeleteConfirmationForm from '../components/dashboard/GenericDeleteConfirmationForm.vue'

const TableProps = Vue.extend({
    props: {
        tableHeaders: Array,
        tableItems: Array,
        resourceName: {
            type: String,
            default: "resources"
        },
        showExpand: {
            type: Boolean,
            default: false,
        },
    }
})

@Component({
    components: {
        GenericDeleteConfirmationForm
    }
})
export default class BaseResourceTable extends mixins(ResourceTableProps, TableProps) {
    selected: any[] = []
    itemToDelete: any | null = null

    get showHideDelete() : boolean {
        return !!this.itemToDelete
    }

    get selectedSet() : Set<any> {
        return new Set<any>(this.selected)
    }

    changeInput(items: any[]) {
        this.$emit('input', items)
    }

    get finalTableHeaders() : any[] {
        let headers = [...this.tableHeaders]

        if (this.useCrudDelete) {
            headers.push({
                text: 'Actions',
                value: 'action',
                sortable: false,
                filterable: false
            })
        }

        if (this.showExpand) {
            headers.push({
                text: '',
                value: 'data-table-expand'
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
        if (this.confirmDelete) {
            this.itemToDelete = item
        } else {
            this.$emit('delete', item)
        }
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
        let children : VNode[] = []
        if (!!this.itemToDelete) {
            children.push(this.$createElement(
                VDialog,
                {
                    props: {
                        persistent: true,
                        "max-width": "40%",
                        value: this.showHideDelete
                    },
                },
                [
                   this.$createElement( 
                        GenericDeleteConfirmationForm,
                        {
                            props: {
                                itemName: this.resourceName,
                                itemsToDelete: [this.itemToDelete.name],
                                useGlobalDeletion: false
                            },
                            on: {
                                'do-cancel': () => { this.itemToDelete = null },
                                'do-delete': () => {
                                    this.$emit('delete', this.itemToDelete)
                                    this.itemToDelete = null
                                }
                            }
                        }
                    )
                ]
            ))
        }

        children.push(this.$createElement(
            VDataTable,
            {
                props: {
                    value: this.selected,
                    headers: this.finalTableHeaders,
                    items: this.tableItems,
                    showSelect: this.selectable,
                    singleSelect: !this.multi,
                    search: this.search,
                    showExpand: this.showExpand,
                },
                on: {
                    input: this.changeInput,
                    'click:row': this.clickRow
                },
                scopedSlots: {
                    ...this.$scopedSlots,
                    'item.action': this.renderActionSlot
                }
            },
        ))

        return this.$createElement('div', children)
    }
}

</script>

<style scoped>

>>>tr {
    cursor: pointer !important;
}

</style>
