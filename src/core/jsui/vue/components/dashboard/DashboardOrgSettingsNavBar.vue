<template>
    <generic-nav-bar :nav-links="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
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
                    title: 'Users',
                    icon: 'mdi-account-multiple',
                    url: PageParamsStore.state.organization!.Url + 'settings/users',
                    disabled : false,
                    hidden: false
                },
                {
                    title: 'Roles',
                    icon: 'mdi-security',
                    url: PageParamsStore.state.organization!.Url + 'settings/roles',
                    disabled : false,
                    hidden: false
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
    watch: {
        mini() {
            this.$emit('on-size-change')
        }
    }
})

</script>

