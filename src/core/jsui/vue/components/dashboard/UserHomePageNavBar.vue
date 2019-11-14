<template>
    <generic-nav-bar :selected-page="selectedPage" :nav-links="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import LocalSettings from '../../../ts/localSettings'
import { createMyProfileUrl, createMyOrgsUrl } from '../../../ts/url'
import GenericNavBar from '../GenericNavBar.vue'
import { PageParamsStore }  from '../../../ts/pageParams'

export default Vue.extend({
    props : {
        selectedPage : Number
    },
    data : function() {
        return {
            navLinks : [
                {
                    title: 'Organizations',
                    icon: 'mdi-account-group',
                    url: createMyOrgsUrl(PageParamsStore.state.user!.Id),
                },
                {
                    title: 'Profile',
                    icon: 'mdi-account-circle',
                    url: createMyProfileUrl(PageParamsStore.state.user!.Id),
                }
            ],
        }
    },
    computed: {
        mini() : boolean {
            return LocalSettings.state.miniNavBar
        }
    },
    components: {
        GenericNavBar
    },
    mounted() {
        VueSetup.store.dispatch('mountPrimaryNavBar', this)
    }
})

</script>
