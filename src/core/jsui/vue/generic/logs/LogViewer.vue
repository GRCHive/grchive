<template>
    <generic-log-viewer
        :full-height="fullHeight"
        :raw-log="rawLog"
    >
    </generic-log-viewer>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GenericLogViewer from './GenericLogViewer.vue'
import { ManagedCode } from '../../../ts/code'
import {
    getLog, TGetLogOutput, TGetLogInput
} from '../../../ts/api/apiLogs'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'

const Props = Vue.extend({
    props: {
        commit: {
            type: String,
            default: "",
        },
        runId: {
            type: Number,
            default: -1,
        },
        runLog: {
            type: Boolean,
            default: false,
        },
        fullHeight: {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        GenericLogViewer
    }
})
export default class LogViewer extends Props {
    rawLog : string | null = null

    refreshLog() {
        this.rawLog = null
        
        let params : TGetLogInput = {
            orgId: PageParamsStore.state.organization!.Id,
        }

        if (this.runId != -1) {
            params.runId = this.runId
            params.runLog = this.runLog
        } else {
            params.commitHash = this.commit
        }

        getLog(params).then((resp : TGetLogOutput) => {
            this.rawLog = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Please refresh the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshLog()
    }
}

</script>
