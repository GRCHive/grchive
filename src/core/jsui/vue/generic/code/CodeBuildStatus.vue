<template>
    <div>
        <script-build-run-status
            v-if="!!status"
            :success="status.Success"
            :start="status.TimeStart"
            :end="status.TimeEnd"
            :show-time-stamp="showTimeStamp"
        >
        </script-build-run-status>

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
import ScriptBuildRunStatus from './ScriptBuildRunStatus.vue'

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

@Component({
    components: {
        ScriptBuildRunStatus
    }
})
export default class CodeBuildStatus extends Props {
    status : DroneCiStatus | null = null

    @Watch('commit')
    refreshStatus() {
        getCodeBuildStatus({
            orgId: PageParamsStore.state.organization!.Id,
            commitHash: this.commit,
        }).then((resp : TGetCodeBuildStatusOutput) => {
            this.status = resp.data

            if (!this.status.Success && !this.status.TimeEnd) {
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
