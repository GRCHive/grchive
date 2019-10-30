<template>
    <generic-nav-bar :selectedPage="selectedPage" :navLinks="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import LocalSettings from '../../../ts/localSettings'
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
                    disabled : false,
                    hidden: true
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
                    url: this.$root.orgUrl + '/risks',
                    disabled : false
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
            return LocalSettings.state.miniNavBar
        }
    },
    mounted() {
        VueSetup.store.dispatch('mountPrimaryNavBar', this)
    },
    watch: {
        mini() {
            this.$emit('on-size-change')
        }
    }
})

</script>
