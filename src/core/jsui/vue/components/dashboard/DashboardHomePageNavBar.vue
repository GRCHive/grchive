<template>
    <generic-nav-bar :selected-page="selectedPage" :nav-links="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import LocalSettings from '../../../ts/localSettings'
import GenericNavBar from '../GenericNavBar.vue'
import { PageParamsStore } from '../../../ts/pageParams'

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
                    url: PageParamsStore.state.organization!.Url,
                    disabled : false,
                    hidden: true
                },
                {
                    title: 'Process Flows',
                    icon: 'mdi-graph-outline',
                    url: PageParamsStore.state.organization!.Url + '/flows',
                    disabled : false
                },
                {
                    title: 'Risks',
                    icon: 'mdi-fire',
                    url: PageParamsStore.state.organization!.Url + '/risks',
                    disabled : false
                },
                {
                    title: 'Controls',
                    icon: 'mdi-shield-lock-outline',
                    url: PageParamsStore.state.organization!.Url + '/controls',
                    disabled : false
                },
                {
                    title: 'General Ledger',
                    icon: 'mdi-bank-outline',
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Access Requests',
                    icon: 'mdi-key-outline',
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Data Sources',
                    icon: 'mdi-database',
                    url: '#',
                    disabled : true
                },
                {
                    title: 'Settings',
                    icon: 'mdi-settings',
                    url: PageParamsStore.state.organization!.Url + '/settings',
                    disabled : false
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
