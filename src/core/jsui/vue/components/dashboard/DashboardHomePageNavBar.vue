<template>
    <generic-nav-bar :selectedPage="selectedPage" :navLinks="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import GenericNavBar from '../GenericNavBar.vue'

export default Vue.extend({
    props : {
        mini : Boolean,
        selectedPage : Number
    },
    data : function() {
        return {
            navLinks : [
                {
                    title: 'Process Flows',
                    icon: 'mdi-graph-outline',
                    //@ts-ignore
                    url: this.$root.orgUrl
                },
                {
                    title: 'Risks',
                    icon: 'mdi-fire',
                    //@ts-ignore
                    url: '#'
                },
                {
                    title: 'Controls',
                    icon: 'mdi-shield-lock-outline',
                    //@ts-ignore
                    url: '#'
                },
                {
                    title: 'General Ledger',
                    icon: 'mdi-bank-outline',
                    //@ts-ignore
                    url: '#'
                },
                {
                    title: 'Access Requests',
                    icon: 'mdi-key-outline',
                    //@ts-ignore
                    url: '#'
                },
                {
                    title: 'Data Sources',
                    icon: 'mdi-database',
                    //@ts-ignore
                    url: '#'
                },
                {
                    title: 'Settings',
                    icon: 'mdi-settings',
                    //@ts-ignore
                    url: '#'
                },
            ],
        }
    },
    components: {
        GenericNavBar
    }
})

</script>
