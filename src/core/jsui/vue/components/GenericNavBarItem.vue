<template>

<v-list-item :href="item.url"
             :to="item.path"
             link
             :color="item.disabled ? `secondary` : primaryColor"
             :disabled="item.disabled"
             :two-line="item.disabled"
             :input-value="isActive"
             v-if="!item.hidden"
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

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { NavBarItem } from '../../ts/navBar'

const VueComponent  = Vue.extend({
    props: {
        item : {
            type: Object as () => NavBarItem,
            default: {} as NavBarItem
        },
        primaryColor: {
            type: String,
            default: "primary"
        },
    }
})

@Component
export default class GenericNavBarItem extends VueComponent {
    get isActive() : boolean {
        if (!this.item.url) {
            return false
        }
        return window.location.pathname.includes(this.item.url)
    }
}

</script>
