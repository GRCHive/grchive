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
            <v-list-item dense>
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
                class="py-0"
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

            <v-btn block text color="primary" class="white-bg">
                See All
            </v-btn>
        </div>
    </v-menu>

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

    markAllNotificationsAsRead() { 
        markNotificationRead({
            userId: PageParamsStore.state.user!.Id,
            notificationIds: this.allNotifications
                .filter((ele : NotificationWrapper) => !ele.Read)
                .map((ele : NotificationWrapper) => ele.Notification.Id)
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

.menu-container {
    min-width: 400px;
    max-width: 400px;
}

.white-bg {
    background-color: white;
}

.notification-label {
    flex: 6 1 !important;
}

</style>
