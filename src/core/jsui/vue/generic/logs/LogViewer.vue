<template>
    <div ref="parent" :style="parentContainerStyle">
        <div v-if="!!rawLog" id="logContainer">
            <div v-html="rawHtml"></div>
        </div>

        <v-progress-circular
            indeterminate
            size="64"
            v-else
        ></v-progress-circular>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { ManagedCode } from '../../../ts/code'
import {
    getLog, TGetLogOutput, TGetLogInput
} from '../../../ts/api/apiLogs'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import Ansi2Html from 'ansi-to-html'

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

@Component
export default class LogViewer extends Props {
    rawLog : string | null = null
    parentBb : DOMRect | null = null

    $refs! : {
        parent : HTMLElement
    }

    get rawHtml() : string {
        if (!this.rawLog) {
            return ""
        }

        const filteredLog = this.rawLog.replace(/\r/g, '\n')
        const converter = new Ansi2Html({
            newline: true,
        })
        return converter.toHtml(filteredLog)
    }

    get parentContainerStyle() : any {
        if (!this.parentBb) {
            Vue.nextTick(() => {
                this.parentBb = <DOMRect>this.$refs.parent.getBoundingClientRect()
            })

            return {}
        }

        let ht = `calc(100vh - ${this.parentBb.top}px)`
        return {
            'height': ht,
            'max-height': ht,
            'overflow': 'auto',
        }
    }

    @Watch('code')
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

<style scoped>

#logContainer {
    overflow: auto;
    font-family: monospace;
    font-size: 12px;
    color: white;
    background-color: #111;
}

</style>
