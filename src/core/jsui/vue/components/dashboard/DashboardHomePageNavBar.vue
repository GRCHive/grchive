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
                    title: 'Business',
                    icon : 'mdi-domain',
                    disabled: false,
                    children: [
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
                            title: 'General Ledger',
                            icon: 'mdi-bank-outline',
                            url: PageParamsStore.state.organization!.Url + 'gl',
                            disabled : false
                        },
                        {
                            title: 'Vendors',
                            icon: 'mdi-store',
                            url: PageParamsStore.state.organization!.Url + 'vendors',
                            disabled : false
                        },
                    ]
                },
                {
                    title: 'GRC',
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
                    ],
                },
                {
                    title: 'Audit',
                    icon: 'mdi-handshake',
                    disabled : false,
                    children: [
                        {
                            title: 'Requests',
                            icon: 'mdi-shield-search',
                            url: PageParamsStore.state.organization!.Url + 'requests',
                            disabled : false
                        },
                        {
                            title: 'Audit Trail',
                            icon: 'mdi-shoe-print',
                            url: PageParamsStore.state.organization!.Url + 'auditlog',
                            disabled : false
                        },
                    ],
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
                        {
                            title: 'Shell',
                            icon: 'mdi-powershell',
                            url: PageParamsStore.state.organization!.Url + 'it/shell',
                        },
                    ]
                },
                {
                    title: 'Automation',
                    icon: 'mdi-cogs',
                    disabled: false,
                    children: [
                        {
                            title: 'Data',
                            icon: 'mdi-table',
                            url: PageParamsStore.state.organization!.Url + 'auto/data',
                        },
                        {
                            title: 'Scripts',
                            icon: 'mdi-code-tags',
                            url: PageParamsStore.state.organization!.Url + 'auto/scripts',
                        },
                        {
                            title: 'Schedule',
                            icon: 'mdi-alarm',
                            url: PageParamsStore.state.organization!.Url + 'auto/schedule',
                        },
                        {
                            title: 'Logs',
                            icon: 'mdi-text-subject',
                            url: PageParamsStore.state.organization!.Url + 'auto/logs',
                        },
                    ],
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
