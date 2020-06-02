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
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import Ansi2Html from 'ansi-to-html'

const Props = Vue.extend({
    props: {
        rawLog: String,
        fullHeight: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class LogViewer extends Props {
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
}

</script>

<style scoped>

#logContainer {
    overflow: auto;
    font-family: monospace;
    font-size: 12px;
    color: white;
    background-color: #111;
    height: 100%;
}

</style>
