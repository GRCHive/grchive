<template>
    <generic-nav-bar :selectedPage="selectedPage" :navLinks="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import GenericNavBar from '../GenericNavBar.vue'

export default Vue.extend({
    props : {
        selectedPage : Number
    },
    data : function() {
        return {
            navLinks : [
                {
                    title: 'Dashboard',
                    icon: 'mdi-view-dashboard',
                    //@ts-ignore
                    url: this.$root.orgUrl,
                    disabled : false
                },
                {
                    title: 'Process Flows',
                    icon: 'mdi-graph-outline',
                    //@ts-ignore
                    url: this.$root.orgUrl + '/flows',
                    disabled : false
                },
                {
                    title: 'Risks',
                    icon: 'mdi-fire',
                    //@ts-ignore
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Controls',
                    icon: 'mdi-shield-lock-outline',
                    //@ts-ignore
                    url: '#',
                    disabled : true
                },
                {
                    title: 'General Ledger',
                    icon: 'mdi-bank-outline',
                    //@ts-ignore
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Access Requests',
                    icon: 'mdi-key-outline',
                    //@ts-ignore
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Data Sources',
                    icon: 'mdi-database',
                    //@ts-ignore
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Settings',
                    icon: 'mdi-settings',
                    //@ts-ignore
                    url: '#',
                    disabled : true
                },
            ],
        }
    },
    components: {
        GenericNavBar
    },
    computed: {
        mini() : boolean {
            return VueSetup.store.state.miniMainNavBar
        }
    },
    mounted() {
        VueSetup.store.dispatch('mountPrimaryNavBar', this)
    }
})

</script>
