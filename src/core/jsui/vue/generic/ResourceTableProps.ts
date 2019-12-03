import Vue from 'vue'

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

export default ResourceTableProps
