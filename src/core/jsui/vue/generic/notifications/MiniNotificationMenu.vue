<template>
    <v-menu
        :close-on-content-click="false"
        offset-y
        bottom
        left
    >
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

        <div class="menu-container">
            <notification-viewer
                :limit="25"
            ></notification-viewer>

            <v-btn
                block
                text
                color="primary"
                class="white-bg"
                :href="allNotificationUrl"
            >
                See All
            </v-btn>
        </div>
    </v-menu>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import NotificationStore from '../../../ts/notifications'
import NotificationViewer from './NotificationViewer.vue'
import { createMyNotificationsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

@Component({
    components: {
        NotificationViewer
    }
})
export default class MiniNotificationMenu extends Vue {
    get hasUnreadNotifications() : boolean {
        return NotificationStore.getters.hasUnreadNotifications
    }

    get allNotificationUrl() : string {
        return createMyNotificationsUrl(PageParamsStore.state.user!.Id)
    }

    mounted() {
        // Probably the safest place to do this since this will be on every page...
        NotificationStore.dispatch('initialize', {
            host: PageParamsStore.state.site!.Host,
        })
    }
}

</script>

<style scoped>

.menu-container {
    min-width: 400px;
    max-width: 400px;
}

.white-bg {
    background-color: white;
}

</style>
