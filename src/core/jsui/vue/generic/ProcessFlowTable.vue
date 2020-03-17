<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { contactUsUrl, createFlowUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'

@Component({
    components: {
        BaseResourceTable,
    }
})
export default class ProcessFlowTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Created',
                value: 'createdTime',
                sort: (a : Date, b : Date) => {
                    return a.getTime() - b.getTime()
                }
            },
            {
                text: 'Last Updated',
                value: 'lastUpdatedTime',
                sort: (a : Date, b : Date) => {
                    return a.getTime() - b.getTime()
                }
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToFlow(item : any) {
        window.location.assign(createFlowUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: inp.Name,
            description: inp.Description,
            createdTime: inp.CreationTime,
            lastUpdatedTime: inp.LastUpdatedTime,
            value: inp
        }
    }

    renderExpansion(props : any) : VNode {
        return this.$createElement(
            'td',
            {
                attrs: {
                    colspan: props.headers.length
                },
            },
            props.item.description
        )
    }

    renderCreatedTime(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.createdTime)
        )
    }

    renderUpdatedTime(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.lastUpdatedTime)
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
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToFlow
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.createdTime': this.renderCreatedTime,
                    'item.lastUpdatedTime': this.renderUpdatedTime,
                }
            }
        )
    }
}

</script>


