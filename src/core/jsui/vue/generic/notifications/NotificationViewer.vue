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

            <!--
            <v-list-item-action>
                <v-btn x-small color="primary" text>
                    Settings
                </v-btn>
            </v-list-item-action>
            -->
        </v-list-item>

        <v-divider></v-divider>
        <v-list
            class="py-0"
            :id="uniqueScrollerId"
            :style="notificationViewerStyle"
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

            <v-list-item
                v-if="loadingMoreNotifications"
            >
                <v-list-item-content>
                    <span class="subtitle-1 flex-center">
                        <v-progress-circular class="center-y" indeterminate size="16"></v-progress-circular>
                        <span class="ml-2">Loading...</span>
                    </span>
                </v-list-item-content>
            </v-list-item>
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

let counter = 0

const Props = Vue.extend({
    props: {
        limit: {
            type: Number,
            default: -1
        },
        useWindowScroll: {
            type: Boolean,
            default: false
        }
    }
})

@Component({
    components: {
        NotificationDisplay
    }
})
export default class NotificationViewer extends Props {
    get allNotifications() : NotificationWrapper[] {
        if (this.limit == -1) {
            return NotificationStore.state.allNotifications
        }
        return NotificationStore.state.allNotifications.slice(0, this.limit)
    }

    get uniqueScrollerId() : string {
        counter += 1
        return `notification-scroll-${counter}`
    }

    get loadingMoreNotifications() : boolean {
        return NotificationStore.state.requestInProgress
    }

    mounted() {
        this.pullMoreNotifications()

        if (this.useWindowScroll) {
            window.addEventListener('scroll', this.handleWheel)
        } else {
            // Not sure why Typescript doesn't like this line?
            //@ts-ignore
            document.querySelector(`#${this.uniqueScrollerId}`)!.addEventListener('scroll', this.handleWheel)
        }
    }

    pullMoreNotifications() {
        if (this.allNotifications.length >= this.limit && this.limit != -1) {
            return
        }

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

    get notificationViewerStyle() : any {
        if (this.useWindowScroll) {
            return {}
        }

        return {
            "max-height": "80vh",
            "overflow": "auto"
        }
    }

    handleWheel(e : Event) {
        let currentScroll: number = 0
        let maxScroll: number = 0

        //@ts-ignore
        let ele : HTMLElement = document.querySelector(`#${this.uniqueScrollerId}`)!

        if (this.useWindowScroll) {
            currentScroll = window.pageYOffset
            maxScroll = document.documentElement.scrollHeight - document.documentElement.clientHeight
        } else {
            currentScroll = ele.scrollTop
            maxScroll = ele.scrollHeight - ele.offsetHeight
        }

        // Just in case we have multiple notification viewers active,
        // we don't really want to trigger multiple pull requests.

        //@ts-ignore
        if (!ele.contains(e.target)) {
            return
        }

        if (currentScroll >= maxScroll) {
            this.pullMoreNotifications()
        }
    }
}

</script>

<style scoped>
.notification-label {
    flex: 6 1 !important;
}

.white-bg {
    background-color: white;
}

.flex-center {
    display: flex;
    justify-content: center;
    align-content: center;
}

</style>
