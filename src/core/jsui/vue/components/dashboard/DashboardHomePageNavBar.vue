<template>
    <generic-nav-bar :nav-links="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import LocalSettings from '../../../ts/localSettings'
import GenericNavBar from '../GenericNavBar.vue'
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    data : function() {
        return {
            navLinks : [
                {
                    title: 'Dashboard',
                    icon: 'mdi-view-dashboard',
                    url: PageParamsStore.state.organization!.Url + 'dashboard',
                    disabled : false,
                    hidden: true
                },
                {
                    title: 'Process Flows',
                    icon: 'mdi-graph-outline',
                    url: PageParamsStore.state.organization!.Url + 'flows',
                    disabled : false
                },
                {
                    title: 'Audit',
                    icon: 'mdi-lock',
                    disabled : false,
                    children: [
                        {
                            title: 'Risks',
                            icon: 'mdi-fire',
                            url: PageParamsStore.state.organization!.Url + 'risks',
                            disabled : false
                        },
                        {
                            title: 'Controls',
                            icon: 'mdi-shield-lock-outline',
                            url: PageParamsStore.state.organization!.Url + 'controls',
                            disabled : false
                        },
                        {
                            title: 'Documentation',
                            icon: 'mdi-file-document-outline',
                            url: PageParamsStore.state.organization!.Url + 'documentation',
                            disabled : false
                        },
                        {
                            title: 'Requests',
                            icon: 'mdi-shield-search',
                            url: PageParamsStore.state.organization!.Url + 'requests',
                            disabled : false
                        },
                    ],
                },
                {
                    title: 'General Ledger',
                    icon: 'mdi-bank-outline',
                    url: PageParamsStore.state.organization!.Url + 'gl',
                    disabled : false
                },
                {
                    title: 'IT',
                    icon: 'mdi-web',
                    disabled : false,
                    children: [
                        {
                            title: 'Systems',
                            icon: 'mdi-application',
                            url: PageParamsStore.state.organization!.Url + 'it/systems',
                        },
                        {
                            title: 'Databases',
                            icon: 'mdi-database',
                            url: PageParamsStore.state.organization!.Url + 'it/databases',
                        },
                        {
                            title: 'Servers',
                            icon: 'mdi-server-network',
                            url: PageParamsStore.state.organization!.Url + 'it/servers',
                        },
                    ]
                },
                {
                    title: 'Settings',
                    icon: 'mdi-settings',
                    url: PageParamsStore.state.organization!.Url + 'settings',
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
