<template>
    <v-card>
        <v-card-title>
            Approval
        </v-card-title>
        <v-divider></v-divider>

        <div class="ma-4 pb-4">
            <p>
                <span class="font-weight-bold">
                    Status:
                </span>
                {{ status }}
                <v-icon small :color="statusIconColor">
                    {{ statusIcon }}
                </v-icon>
            </p>

            <div v-if="!!approval">
                <user-search-form-component
                    label="Responder"
                    v-bind:user="responderUser"
                    readonly
                ></user-search-form-component>

                <p>
                    <span class="font-weight-bold">
                        Responded At:
                    </span>
                    {{ responseTime }}
                </p>

                <p v-if="!approval.Response">
                    <span class="font-weight-bold">
                        Reason:
                    </span>
                    <pre class="pb-4">{{ approval.Reason }}</pre>
                </p>
            </div>
        </div>
    </v-card>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { GenericApproval } from '../../../ts/requests'
import { standardFormatTime } from '../../../ts/time'
import UserSearchFormComponent from '../UserSearchFormComponent.vue'
import MetadataStore from '../../../ts/metadata'

const Props = Vue.extend({
    props : {
        approval: {
            type: Object,
            default: () => null as GenericApproval | null
        }
    }
})

@Component({
    components: {
        UserSearchFormComponent,
    }
})
export default class GenericApprovalDisplay extends Props {
    get status() : string {
        if (!this.approval) {
            return "Pending"
        }

        if (!this.approval.Response) {
            return "Denied"
        } else {
            return "Approved"
        }
    }

    get statusIcon() : string {
        if (!this.approval) {
            return "mdi-help-circle"
        }

        if (!this.approval.Response) {
            return "mdi-cancel"
        } else {
            return "mdi-check"
        }
    }

    get statusIconColor() : string {
        if (!this.approval) {
            return "warning"
        }

        if (!this.approval.Response) {
            return "error"
        } else {
            return "success"
        }
    }

    get responderUser() : User | null {
        if (!this.approval) {
            return null
        }
        return MetadataStore.getters.getUser(this.approval!.ResponderUserId)
    }

    get responseTime() : string {
        if (!this.approval) {
            return "Unknown"
        }
        return standardFormatTime(this.approval!.ResponseTime)
    }
}

</script>
