<template>
    <v-navigation-drawer app clipped
        :mini-variant="mini"
        mini-variant-width=50
        :width="width"
        permanent
    >
        <slot></slot>
        <v-list class="py-0">
            <v-list-item-group :value="selectedPage" mandatory>
                <v-list-item v-for="(item, i) in navLinks" 
                             :key="i"
                             :href="item.url"
                             :to="item.path"
                             link
                             :color="item.disabled ? `secondary` : primaryColor"
                             :disabled="item.disabled"
                             :two-line="item.disabled"
                             @click="doItemClick($event, i)"
                >
                    <v-list-item-icon v-if="item.icon != ''">
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

                    <v-list-item-action v-if="!!item.action">
                        <v-btn icon @click.stop="item.action.fn" @mousedown.stop>
                            <v-icon>{{ item.action.icon }}</v-icon>
                        </v-btn>
                    </v-list-item-action>
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
        navLinks: {},
        primaryColor: {
            type: String,
            default: "primary"
        },
        width: {
            type: Number,
            default: 256
        }
    },
    methods: {
        doItemClick(e : MouseEvent, idx : number) {
            this.$emit('item-change', e, idx)
        }
    }
})

</script>
