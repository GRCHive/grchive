<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const ResourceTableProps = Vue.extend({
    props: {
        resources: Array,
        value: {
            type: Array,
            default: () => []
        },
        selectable: {
            type: Boolean,
            default: false
        },
        multi: {
            type: Boolean,
            default: false
        },
        search : {
            type: String,
            default: ""
        }
    }
})

@Component
export default class BaseResourceTable extends ResourceTableProps {
    selected = []

    get tableHeaders() : any[] {
        return []
    }

    get tableItems() : any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    changeInput(items: any[]) {
        this.$emit('input', items.map(this.transformTableItemToInputResource))
    }

    transformInputResourceToTableItem(inp : any) : any {
        return null
    }

    transformTableItemToInputResource(inp : any) : any {
        return null
    }
}

</script>
