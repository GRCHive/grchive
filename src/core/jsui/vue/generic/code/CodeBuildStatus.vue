<template>
    <div>
        <div class="statusContainer" v-if="!!status">
            <v-icon
                :color="iconColor"
                small
            >
                {{ iconName }}
            </v-icon>

            <div v-if="showTimeStamp" class="ml-2">
                <p class="ma-0">
                    <span class="font-weight-bold">Start: </span>
                    {{ timeStartStr }}
                </p>

                <p class="ma-0" v-if="!isPending">
                    <span class="font-weight-bold">End: </span>
                    {{ timeEndStr }}
                </p>
            </div>
        </div>

        <v-progress-circular
            indeterminate
            size="16"
            v-else
        ></v-progress-circular>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { DroneCiStatus } from '../../../ts/code'
import {
    getCodeBuildStatus, TGetCodeBuildStatusOutput,
} from '../../../ts/api/apiCode'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { standardFormatTime } from '../../../ts/time'

const Props = Vue.extend({
    props: {
        commit: String,
        showTimeStamp: {
            type: Boolean,
            default: false,
        }
    }
})

const refreshPeriodMs = 5000

@Component
export default class CodeBuildStatus extends Props {
    status : DroneCiStatus | null = null

    get isPending() {
        return this.status!.TimeEnd == null
    }

    get isSuccess() {
        return this.status!.Success
    }

    get timeStartStr() : string {
        return standardFormatTime(this.status!.TimeStart)
    }

    get timeEndStr() : string {
        return standardFormatTime(this.status!.TimeEnd!)
    }

    get iconColor() {
        if (this.isPending) {
            return 'warning'
        } else if (this.isSuccess) {
            return 'success'
        } else {
            return 'error'
        }
    }

    get iconName() {
        if (this.isPending) {
            return 'mdi-help-circle'
        } else if (this.isSuccess) {
            return 'mdi-check-circle'
        } else {
            return 'mdi-close-circle'
        }
    }

    @Watch('commit')
    refreshStatus() {
        getCodeBuildStatus({
            orgId: PageParamsStore.state.organization!.Id,
            commitHash: this.commit,
        }).then((resp : TGetCodeBuildStatusOutput) => {
            this.status = resp.data

            if (this.isPending) {
                setTimeout(this.refreshStatus, refreshPeriodMs)
            }
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

    mounted() {
        this.refreshStatus()
    }
}

</script>

<style scoped>

.statusContainer {
    display: flex;
    align-items: center;
    justify-items: center;
}

</style>
