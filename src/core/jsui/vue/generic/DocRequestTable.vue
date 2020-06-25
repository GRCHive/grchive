<script lang="ts">

import Vue, {VNode} from 'vue'
import { VIcon, VChip } from 'vuetify/lib'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { createSingleDocRequestUrl, contactUsUrl } from '../../ts/url'
import { TGetDocCatOutput, getDocumentCategory } from '../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import DocRequestStatusDisplay from './requests/DocRequestStatusDisplay.vue'

@Component({
    components: {
        BaseResourceTable,
        DocRequestStatusDisplay
    }
})
export default class DocRequestTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Requester',
                value: 'requester',
            },
            {
                text: 'Request Time',
                value: 'requestTime',
                sort: (a : Date, b : Date) => {
                    return a.getTime() - b.getTime()
                }
            },
            {
                text: 'Assignee',
                value: 'assignee',
            },
            {
                text: 'Due Date',
                value: 'dueDate',
                sort: (a : Date, b : Date) => {
                    return a.getTime() - b.getTime()
                }
            },
            {
                text: 'Status',
                value: 'status',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToDocRequest(item : any) {
        window.location.assign(createSingleDocRequestUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: inp.Name,
            requester: createUserString(MetadataStore.getters.getUser(inp.RequestedUserId)),
            requestTime: inp.RequestTime,
            assignee: createUserString(MetadataStore.getters.getUser(inp.AssigneeUserId)),
            dueDate: inp.DueDate,
            value: inp
        }
    }

    renderStatus(props : any) : VNode {
        return this.$createElement(
            DocRequestStatusDisplay,
            {
                props: {
                    documentRequest: props.item.value,
                }
            },
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

    renderRequestTime(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.requestTime)
        )
    }

    renderDueDate(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.dueDate)
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
                    'click:row': this.goToDocRequest
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.status': this.renderStatus,
                    'item.requestTime': this.renderRequestTime,
                    'item.dueDate': this.renderDueDate,
                }
            }
        )
    }
}

</script>

