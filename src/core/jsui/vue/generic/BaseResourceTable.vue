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
    selected: any[] = []

    get valueSet() : Set<any> {
        return new Set<any>(this.value)
    }

    get tableHeaders() : any[] {
        return []
    }

    get tableItems() : any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    changeInput(items: any[]) {
        this.$emit('input', Array.from(new Set(items)).map(this.transformTableItemToInputResource))
    }

    manualToggleItem(item : any) {
        let val = this.transformTableItemToInputResource(item)
        let newValueArr = []
        if (this.valueSet.has(val)) {
            newValueArr = this.value.filter((ele : any) => ele != val)
        } else {
            newValueArr = [...this.value, val]
        }

        this.$emit('input', newValueArr)
        this.selected = newValueArr.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return null
    }

    transformTableItemToInputResource(inp : any) : any {
        return null
    }

}

</script>
