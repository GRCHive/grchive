<template>
    <div>
        <notification-popup
            v-for="(item, index) in recentNotifications"
            :key="index"
            :notification="item"
            :style="styles[index]"
            :ref="`popup${index}`"
            @onready="recomputeHeightForNotification(index)"
        >
        </notification-popup>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import NotificationPopup from './NotificationPopup.vue'
import NotificationStore, { NotificationWrapper } from '../../../ts/notifications'

@Component({
    components: {
        NotificationPopup
    }
})
export default class NotificationPopupManager extends Vue {
    notificationHeights: number[] = []

    @Watch('recentNotifications')
    recomputeNotificationHeights() {
        this.notificationHeights = this.recentNotifications.map(() => 80)

        for (let i = 0; i < this.notificationHeights.length; ++i) {
            this.recomputeHeightForNotification(i)
        }
    }

    recomputeHeightForNotification(i : number) {
        Vue.nextTick(() => {
            //@ts-ignore
            let ele : HTMLElement = this.$refs[`popup${i}`][0].$el
            Vue.set(this.notificationHeights, i, ele.offsetHeight)
        })
    }

    get recentNotifications() : NotificationWrapper[] {
        // 5 is probably a good number of notifications to show...
        return NotificationStore.state.recentNotifications.slice(0, 5)
    }

    get styles() : any[] {
        const margin : number = 8
        let ret : any[] = []
        let offsetY : number = 12
        for (let i = 0; i < this.recentNotifications.length; ++i) {
            offsetY += this.notificationHeights[i] + margin

            ret.push({
                "top": `calc(100vh - ${offsetY}px)`,
                "right": "20px",
            })
        }
        return ret
    }

    mounted() {
        this.recomputeNotificationHeights()
    }
}

</script>
