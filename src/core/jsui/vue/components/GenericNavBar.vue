<template>
    <v-navigation-drawer app clipped
        :mini-variant.sync="mini"
        mini-variant-width=50
    >
        <v-list class="py-0">
            <v-list-item-group :value="selectedPage">
                <v-list-item v-for="(item, i) in navLinks" 
                             :key="i"
                             :href="item.url"
                             link
                             :color="item.disabled ? `secondary` : `primary`"
                             :disabled="item.disabled"
                             :two-line="item.disabled"
                >
                    <v-list-item-icon>
                        <v-icon>{{ item.icon }}</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <v-list-item-title>
                            {{ item.title }}
                        </v-list-item-title>
                        <v-list-item-subtitle v-if="item.disabled">
                            Coming Soon.
                        </v-list-item-subtitle>
                    </v-list-item-content>
                </v-list-item>
            </v-list-item-group>
        </v-list>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
    props : {
        mini : Boolean,
        selectedPage : Number,
        navLinks: {}
    },
})

</script>
