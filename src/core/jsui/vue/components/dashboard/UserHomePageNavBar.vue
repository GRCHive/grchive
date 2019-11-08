<template>
    <generic-nav-bar :selectedPage="selectedPage" :navLinks="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import LocalSettings from '../../../ts/localSettings'
import { createMyAccountUrl } from '../../../ts/url'
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
                    title: 'Profile',
                    icon: 'mdi-account-circle',
                    url: createMyAccountUrl(PageParamsStore.state.user!.Email),
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
