<template>
    <v-navigation-drawer class="max-height">
        <v-list class="py-0">
            <v-list-item-group :value="selectedPage">
                <v-list-item v-for="(item, i) in navLinks" 
                             :key="i"
                             :href="item.url"
                             link
                >
                    <v-list-item-icon>
                        <v-icon>{{ item.icon }}</v-icon>
                    </v-list-item-icon>

                    <v-list-item-title>
                        {{ item.title }}
                    </v-list-item-title>
                </v-list-item>
            </v-list-item-group>
        </v-list>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import { createMyAccountUrl } from '../../../ts/url'

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
                    //@ts-ignore
                    url: createMyAccountUrl(this.$root.userEmail),
                }
            ],
        }
    }
})

</script>
