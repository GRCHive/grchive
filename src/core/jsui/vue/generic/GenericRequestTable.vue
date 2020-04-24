<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { 
    contactUsUrl,
} from '../../ts/url'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { GenericApproval } from '../../ts/requests'
import { VIcon, VTooltip } from 'vuetify/lib'
import { standardFormatTime } from '../../ts/time'
import { getGenericApproval, TGetApprovalOutput } from '../../ts/api/apiRequests'
import { PageParamsStore } from '../../ts/pageParams'

@Component({
    components: {
        BaseResourceTable,
    }
})
export default class GenericRequestTable extends ResourceTableProps {
    idToApproval : Record<number, GenericApproval> = Object()
    pendingApprovalRequests : Set<number> = new Set<number>()

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
                sort: (a : Date | null, b : Date | null) => {
                    if (!a) { 
                        return 1
                    }

                    if (!b) {
                        return -1
                    }

                    return a.getTime() - b.getTime()
                }
            },
            {
                text: 'Approval',
                value: 'approval',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    loadApproval(id : number) {
        if (this.pendingApprovalRequests.has(id)) {
            return
        }

        this.pendingApprovalRequests.add(id)
        getGenericApproval({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: id,
        }).then((resp : TGetApprovalOutput) => {
            Vue.set(this.idToApproval, id, resp.data)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    transformInputResourceToTableItem(inp : any) : any {
        this.loadApproval(inp.Id)

        return {
            id: inp.Id,
            name: inp.Name,
            requester: createUserString(MetadataStore.getters.getUser(inp.UploadUserId)),
            requestTime: inp.UploadTime,
            assignee: createUserString(MetadataStore.getters.getUser(inp.Assignee)),
            dueDate: inp.DueDate,
            approval: this.idToApproval[inp.Id],
            value: inp,
        }
    }

    renderRequestTime(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.requestTime)
        )
    }

    renderApproval(props : any) : VNode {
        let approval : GenericApproval | null = props.item.approval

        if (!approval) {
            return this.$createElement(
                VIcon,
                {
                    props: {
                        small : true,
                        color: 'warning'
                    }
                },
                'mdi-help-circle'
            )
        } else if (approval.Response) {
            return this.$createElement(
                VIcon,
                {
                    props: {
                        small : true,
                        color: 'success'
                    }
                },
                'mdi-check'
            )
        } else {
            let renderIcon = (props : any) => {
                return this.$createElement(
                    VIcon,
                    {
                        props: {
                            small : true,
                            color: 'error'
                        },
                        on: props.on,
                    },
                    'mdi-close'
                )
            }

            return this.$createElement(
                VTooltip,
                {
                    props: {
                        bottom: true,
                    },
                    scopedSlots: {
                        activator: renderIcon
                    }
                },
                approval.Reason
            )
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
            props.item.value.Description
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
                    resourceName: "generic request",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                },
                scopedSlots: {
                    'item.approval': this.renderApproval,
                    'item.requestTime': this.renderRequestTime,
                    'item.dueDate': this.renderDueDate,
                    'expanded-item': this.renderExpansion,
                }
            }
        )
    }

}

</script>
