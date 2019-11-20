<template>
    <v-navigation-drawer app clipped
        :mini-variant="mini"
        mini-variant-width=50
        :width="width"
        permanent
    >
        <slot></slot>
        <v-list class="py-0">
            <div v-for="(item, i) in navLinks" 
                 :key="i"
                 :style="!!item.hidden ? `display: none;` : ``"
            >
                <v-list-group v-if="!!item.children && item.children.length > 0"
                              :prepend-icon="item.icon"
                              no-action
                              value="true"
                              color="rgba(0, 0, 0, 0.87) !important"
                >

                    <template v-slot:activator>
                        <v-list-item-title>{{ item.title }} </v-list-item-title>
                    </template>

                    <generic-nav-bar-item v-for="(child, ci) in item.children"
                                          :key="ci"
                                          :item="child"
                                          :primary-color="primaryColor">
                    </generic-nav-bar-item>
                </v-list-group>

                <generic-nav-bar-item :item="item"
                                      :primary-color="primaryColor"
                                      v-else>
                </generic-nav-bar-item>
            </div>
        </v-list>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import GenericNavBarItem from './GenericNavBarItem.vue'

export default Vue.extend({
    props : {
        mini : Boolean,
        kelectedPage : Number,
        navLinks: Array,
        primaryColor: {
            type: String,
            default: "primary"
        },
        width: {
            type: Number,
            default: 256
        }
    },
    components: {
        GenericNavBarItem
    },
})

</script>

<style scoped>

>>>.v-list-group__header .v-list-group__header__prepend-icon .v-icon {
    color: rgba(0, 0, 0, 0.54) !important;
}

</style>
