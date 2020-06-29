<template>
    <v-chip
        small
        :color="colorStr"
        :close="canClose"
        @click:close="onClose"
    >
        {{ statusStr }}
        <v-icon
            small
            right
        >
            {{ iconStr }}
        </v-icon>
    </v-chip>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import { DocRequestStatus, getDocumentRequestStatus } from '../../../ts/docRequests'

@Component
export default class DocRequestStatusDisplay extends Vue {
    @Prop()
    status! : DocRequestStatus

    @Prop({default:false})
    canClose! : boolean

    onClose() {
        this.$emit('click:close')
    }

    get iconStr() : string {
        switch (this.status) {
            case DocRequestStatus.Open:
                return "mdi-circle-outline"
            case DocRequestStatus.InProgress:
                return "mdi-progress-wrench"
            case DocRequestStatus.Feedback:
                return "mdi-message-alert-outline"
            case DocRequestStatus.Complete:
                return "mdi-account-check-outline"
            case DocRequestStatus.Overdue:
                return "mdi-alert-circle-outline"
            case DocRequestStatus.Approved:
                return "mdi-check-circle-outline"
        }

        return ''
    }

    get statusStr() : string {
        switch (this.status) {
            case DocRequestStatus.Open:
                return "Open"
            case DocRequestStatus.InProgress:
                return "In Progress"
            case DocRequestStatus.Feedback:
                return "Feedback"
            case DocRequestStatus.Complete:
                return "Pending Approval"
            case DocRequestStatus.Overdue:
                return "Overdue"
            case DocRequestStatus.Approved:
                return "Approved"
        }

        return ''
    }

    get colorStr() : string {
        switch (this.status) {
            case DocRequestStatus.Open:
                return "secondary"
            case DocRequestStatus.InProgress:
                return "primary"
            case DocRequestStatus.Feedback:
                return "warning"
            case DocRequestStatus.Complete:
                return "accent"
            case DocRequestStatus.Overdue:
                return "error"
            case DocRequestStatus.Approved:
                return "success"
        }
        return ''
    }
}

</script>
