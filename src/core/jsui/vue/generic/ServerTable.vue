<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleServerUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class ServerTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Operating System',
                value: 'os',
            },
            {
                text: 'Location',
                value: 'location',
            },
            {
                text: 'IP Address',
                value: 'ip',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: inp.Name,
            os: inp.OperatingSystem,
            location: inp.Location,
            ip: inp.IpAddress,
            value: inp
        }
    }

    goToServer(item : any) {
        window.location.assign(createSingleServerUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
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
                    'click:row': this.goToServer
                }
            }
        )
    }
}

</script>
