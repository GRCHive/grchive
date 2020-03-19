<template>
    <v-menu offset-y bottom left>
        <template v-slot:activator="{on}">
            <v-badge
                color="red"
                dot
                overlap
                :value="hasUnreadNotifications"
                :offset-x="16"
                :offset-y="16"
            >
                <v-btn color="primary" icon v-on="on">
                    <v-icon>mdi-bell</v-icon>
                </v-btn>
            </v-badge>
        </template>

        <v-list
            dense
            tile
            min-width="400"
            max-width="400"
        >
            <template v-for="(item, index) in allNotifications">
                <notification-display
                    :key="`display-${index}`"
                    :notification="item"
                >
                </notification-display>
                <v-divider :key="`divider-${index}`"></v-divider>
            </template>
        </v-list>
    </v-menu>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import NotificationStore, { NotificationWrapper } from '../../../ts/notifications'
import NotificationDisplay from './NotificationDisplay.vue'

@Component({
    components: {
        NotificationDisplay
    }
})
export default class MiniNotificationMenu extends Vue {

    get allNotifications() : NotificationWrapper[] {
        return NotificationStore.state.allNotifications
    }

    get hasUnreadNotifications() : boolean {
        return NotificationStore.getters.hasUnreadNotifications
    }

    mounted() {
        NotificationStore.dispatch('pullNotifications')
    }
}

</script>
