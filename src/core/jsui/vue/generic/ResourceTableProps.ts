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
        },
        useCrudDelete: {
            type: Boolean,
            default: false
        },
        confirmDelete: {
            type: Boolean,
            default: false
        },
        useGlobalDeletion: {
            type: Boolean,
            default: false
        }
    }
})

export default ResourceTableProps
