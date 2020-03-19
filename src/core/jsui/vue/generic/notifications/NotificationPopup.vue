<template>
    <div class="notification-popup">
        <v-card>
            <notification-display
                :notification="notification"
                force-unread-display
                enable-close
                @close="close"
            >
            </notification-display>
        </v-card>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import NotificationStore, { NotificationWrapper } from '../../../ts/notifications'
import NotificationDisplay from './NotificationDisplay.vue'

const Props = Vue.extend({
    props: {
        notification: {
            type: Object,
            default: () => Object() as NotificationWrapper
        }
    }
})

@Component({
    components: {
        NotificationDisplay
    }
})
export default class NotificationPopup extends Props {
    close() {
        NotificationStore.commit('removeRecentNotification', this.notification.Notification.Id)
    }
}

</script>

<style scoped>

.notification-popup {
    position: fixed;
    max-width: 400px;
    min-width: 400px;
}

</style>
