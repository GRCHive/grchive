<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { standardFormatDate } from '../../ts/time'
import { createSingleDocFileUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class DocFileTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name'
            },
            {
                text: 'Filename',
                value: 'filename'
            },
            {
                text: 'Relevant Time',
                value: 'relevantTime'
            },
            {
                text: 'Upload Time',
                value: 'uploadTime'
            },
            {
                text: 'User',
                value: 'user'
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToDocFile(item : any) {
        window.location.assign(createSingleDocFileUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: inp.AltName,
            filename: inp.StorageName,
            relevantTime: standardFormatDate(inp.RelevantTime),
            uploadTime: standardFormatDate(inp.UploadTime),
            user: createUserString(MetadataStore.getters.getUser(inp.UploadUserId)),
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
                    resourceName: "files",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToDocFile
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion
                }
            }
        )
    }
}

</script>


