<template>
    <div>
        <v-list-item dense class="white-bg">
            <v-list-item-content class="notification-label">
                <v-list-item-title>
                    Notifications
                </v-list-item-title>
            </v-list-item-content>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn
                    x-small
                    color="primary"
                    text
                    @click="markAllNotificationsAsRead"
                >
                    Mark All As Read
                </v-btn>
            </v-list-item-action>

            <v-list-item-action>
                <v-btn x-small color="primary" text>
                    Settings
                </v-btn>
            </v-list-item-action>
        </v-list-item>

        <v-divider></v-divider>
        <v-list
            class="py-0 notif-container"
            dense
            tile
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
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import NotificationStore, { NotificationWrapper } from '../../../ts/notifications'
import NotificationDisplay from './NotificationDisplay.vue'
import { markNotificationRead } from '../../../ts/api/apiNotifications'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

@Component({
    components: {
        NotificationDisplay
    }
})
export default class NotificationViewer extends Vue {
    get allNotifications() : NotificationWrapper[] {
        return NotificationStore.state.allNotifications
    }

    mounted() {
        NotificationStore.dispatch('pullNotifications')
    }

    markAllNotificationsAsRead() { 
        markNotificationRead({
            userId: PageParamsStore.state.user!.Id,
            notificationIds: [],
            all: true,
        }).then(() => {
            NotificationStore.commit('markAllAsRead')
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

}

</script>

<style scoped>
.notification-label {
    flex: 6 1 !important;
}

.notif-container {
    overflow: auto;
}

.white-bg {
    background-color: white;
}


</style>
