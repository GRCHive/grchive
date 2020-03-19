<template>
    <v-menu offset-y bottom left>
        <template v-slot:activator="{on}">
            <v-badge
                color="red"
                bottom
                dot
                overlap
                :value="hasNotifications"
                :offset-x="16"
                :offset-y="16"
            >
                <v-btn color="primary" icon v-on="on">
                    <v-icon>mdi-bell</v-icon>
                </v-btn>
            </v-badge>
        </template>

        <v-list>
        </v-list>
    </v-menu>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

@Component
export default class MiniNotificationMenu extends Vue {

    get hasNotifications() : boolean {
        return true
    }
}

</script>
