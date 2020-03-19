<template>
    <div 
        :class="'d-flex ' + ((notification.Read && !forceUnreadDisplay) ? '' : 'unread-notification')"
    >
        <v-list-item
            :href="resourceUrl"
            @click.stop="onClick"
            v-if="ready"
            two-line
        >
            <v-list-item-icon class="notification-icon">
                <v-icon>
                    {{ resourceIcon }}
                </v-icon>
            </v-list-item-icon>

            <v-list-item-content>
                <v-list-item-title
                    class="long-text"
                    v-text="displayText"
                >
                </v-list-item-title>

                <v-list-item-subtitle>
                    {{ timestampStr }}
                </v-list-item-subtitle>
            </v-list-item-content>
        </v-list-item>

        <v-row align="center" justify="center" class="py-4" v-else>
            <v-progress-circular indeterminate size="16"></v-progress-circular>
        </v-row>

        <v-btn icon @click.stop="close" class="center-y" v-if="enableClose">
            <v-icon small>
                mdi-close-circle-outline
            </v-icon>
        </v-btn>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { NotificationWrapper } from '../../../ts/notifications'
import { ResourceHandle, resourceTypeToIcon } from '../../../ts/resourceUtils'
import { TGetResourceHandleOutput , getResourceHandle } from '../../../ts/api/apiResources'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { standardFormatTime } from '../../../ts/time'
import { markNotificationRead } from '../../../ts/api/apiNotifications'

const Props = Vue.extend({
    props: {
        notification: {
            type: Object,
            default: () => Object() as NotificationWrapper
        },
        forceUnreadDisplay : {
            type: Boolean,
            default: false,
        },
        enableClose: {
            type: Boolean,
            default: false
        }
    }
})

@Component
export default class NotificationDisplay extends Props {
    subjectHandle: ResourceHandle | null = null
    objectHandle : ResourceHandle | null = null
    indirectObjectHandle: ResourceHandle | null = null

    get resourceIcon() : string {
        return resourceTypeToIcon(this.notification.Notification.ObjectType)
    }

    get displayText() : string {
        if (!this.ready) {
            return ""
        }

        let text = [
            !!this.subjectHandle ? this.subjectHandle.displayText : 'Someone',
        ]
        text.push(this.notification.Notification.Verb)

        if (!!this.objectHandle) {
            text.push(this.objectHandle.displayText)
        }

        if (!!this.indirectObjectHandle) {
            text.push('to')
            text.push(this.indirectObjectHandle.displayText)
        }

        return text.join(' ') + '.'
    }

    get resourceUrl() : string {
        return (!!this.objectHandle && !!this.objectHandle.resourceUri) ?
            this.objectHandle.resourceUri :
            '#'
    }

    get ready() : boolean {
        return (!!this.subjectHandle || this.notification.Notification.SubjectType == "") && 
            (!!this.objectHandle || this.notification.Notification.ObjectType == "") &&
            (!!this.indirectObjectHandle || this.notification.Notification.IndirectObjectType == "")
    }

    get timestampStr() : string {
        return standardFormatTime(this.notification.Notification.Time)
    }

    onError() {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops! Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    onClick() {
        markNotificationRead({
            userId: PageParamsStore.state.user!.Id,
            notificationIds: [this.notification.Notification.Id],
            all: false
        })
    }

    refreshResourceHandles() {
        if (this.notification.Notification.SubjectType != "") {
            getResourceHandle({
                orgId: this.notification.Notification.OrgId,
                resourceType: this.notification.Notification.SubjectType,
                resourceId: this.notification.Notification.SubjectId,
            }).then((resp : TGetResourceHandleOutput) => {
                this.subjectHandle = resp.data
            }).catch(this.onError)
        } 

        if (this.notification.Notification.ObjectType != "") {
            getResourceHandle({
                orgId: this.notification.Notification.OrgId,
                resourceType: this.notification.Notification.ObjectType,
                resourceId: this.notification.Notification.ObjectId,
            }).then((resp : TGetResourceHandleOutput) => {
                this.objectHandle = resp.data
            }).catch(this.onError)
        } 

        if (this.notification.Notification.IndirectObjectType != "") {
            getResourceHandle({
                orgId: this.notification.Notification.OrgId,
                resourceType: this.notification.Notification.IndirectObjectType,
                resourceId: this.notification.Notification.IndirectObjectId,
            }).then((resp : TGetResourceHandleOutput) => {
                this.indirectObjectHandle = resp.data
            }).catch(this.onError)
        } 
    }

    mounted() {
        this.refreshResourceHandles()
    }

    close() {
        this.$emit('close')
    }
}

</script>

<style scoped>

.unread-notification {
    background-color: #BBDEFB;
}

.notification-icon {
    margin-top: auto !important;
    margin-bottom: auto !important;
}

</style>
