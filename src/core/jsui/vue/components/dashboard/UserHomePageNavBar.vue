<template>
    <generic-nav-bar :selectedPage="selectedPage" :navLinks="navLinks" :mini="mini">
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import { createMyAccountUrl } from '../../../ts/url'
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
                    title: 'Profile',
                    icon: 'mdi-account-circle',
                    //@ts-ignore
                    url: createMyAccountUrl(this.$root.userEmail),
                }
            ],
        }
    },
    components: {
        GenericNavBar
    }
})

</script>
