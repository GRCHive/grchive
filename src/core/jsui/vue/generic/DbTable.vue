<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import { createSingleDbUrl } from '../../ts/url'
import { getDbTypeAsString } from '../../ts/databases'
import MetadataStore from '../../ts/metadata'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class DbTable extends ResourceTableProps {
    get ready() : boolean {
        return MetadataStore.state.dbTypesInitialized
    }

    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Type',
                value: 'type',
            },
            {
                text: 'Version',
                value: 'version',
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
            type: getDbTypeAsString(inp),
            version: inp.Version,
            value: inp
        }
    }

    goToDb(item : any) {
        window.location.assign(createSingleDbUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }

    render() : VNode {
        if (this.ready) {
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
                        'click:row': this.goToDb
                    }
                }
            )
        } else {
            return this.$createElement('div')
        }
    }
}

</script>
