<template>
    <div>
        <generic-code-toolbar
            @save="onSave"
            ref="toolbar"
            :code-value="value"
            :save-in-progress="saveInProgress"
        >
            <template v-slot:custom-menu>
                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn text color="accent" v-on="on">
                            Logs
                            <v-icon small color="accent">mdi-chevron-down</v-icon>
                        </v-btn>
                    </template>
                    <v-list dense>
                        <v-list-item dense @click="showLogs = !showLogs">
                            <v-list-item-action>
                                <v-checkbox :input-value="showLogs">
                                </v-checkbox>
                            </v-list-item-action>
                            <v-list-item-title>
                                Build Logs
                            </v-list-item-title>
                        </v-list-item>
                        <v-divider></v-divider>
                    </v-list>
                </v-menu>
            </template>
            <template v-slot:custom-status>
                <v-col cols="auto">
                    <v-select
                        dense
                        solo
                        flat
                        hide-details
                        :items="allCodeItems"
                        label="Version"
                        v-model="selectedCode"
                    >
                    </v-select>
                </v-col>
            </template>
        </generic-code-toolbar>
        <v-divider></v-divider>

        <dynamic-split-container :enable-col-b="showLogs">
            <template v-slot:first-col>
                <generic-code-editor
                    :value="codeString"
                    :lang="lang"
                    :readonly="readonly"
                    :full-height="fullHeight"
                    @input="onInput"
                    ref="editor"
                >
                </generic-code-editor>
            </template>

            <template v-slot:second-col>
                <log-viewer
                    :code="selectedCode"
                    full-height
                >
                </log-viewer>
            </template>
        </dynamic-split-container>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component, { mixins } from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import GenericCodeEditor, { Props } from './GenericCodeEditor.vue'
import GenericCodeToolbar from './GenericCodeToolbar.vue'
import { ManagedCode } from '../../../ts/code'
import { PageParamsStore } from '../../../ts/pageParams'
import { 
    getCode, TGetCodeOutput,
    allCode, TAllCodeOutput,
    saveCode, TSaveCodeOutput,
} from '../../../ts/api/apiCode'
import { contactUsUrl } from '../../../ts/url'
import { standardFormatTime } from '../../../ts/time'
import LogViewer from '../logs/LogViewer.vue'
import DynamicSplitContainer from '../DynamicSplitContainer.vue'

const ManagedProps = Vue.extend({
    props: {
        dataId: {
            type: Number,
            default: -1,
        },
        scriptId: {
            type: Number,
            default: -1,
        },
    }
})

@Component({
    components: {
        LogViewer,
        DynamicSplitContainer,
        GenericCodeToolbar,
        GenericCodeEditor,
    }
})
export default class ManagedCodeIDE extends mixins(Props, ManagedProps) {
    codeString : string = ""
    loading: boolean = false
    showLogs: boolean = false

    // This is equivalent to loading for the first time
    // we load code. We need this because the code toolbar
    // will use 'codeString' when mounted to determine if the user
    // needs to save. Thus, we need to delay its mount until after
    // the code loads. We don't want to use 'loading' to prevent
    // the mount because otherwise there's a flicker every time you
    // change the code version which is slightly unpleasant.
    initialLoad: boolean = true

    saveInProgress: boolean = false

    allCode : ManagedCode[] = []
    selectedCode : ManagedCode | null = null

    $refs!: {
        editor: GenericCodeEditor,
        toolbar: GenericCodeToolbar,
    }

    onInput(text : string) {
        this.codeString = text
    }

    get allCodeItems() : any[] {
        return this.allCode.map((ele : ManagedCode, idx : number) => ({
            text: `v${this.allCode.length-idx} [${standardFormatTime(ele.ActionTime)}]`,
            value: ele, 
        }))
    }

    @Watch('selectedCode')
    pullCode() {
        if (!this.selectedCode) {
            this.codeString = ""
            this.initialLoad = false
            return
        }

        this.loading = true
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
            codeId: this.selectedCode.Id,
        }

        if (this.dataId != -1) {
            params.dataId = this.dataId
        }

        if (this.scriptId != -1) {
            params.scriptId = this.scriptId
        }

        getCode(params).then((resp : TGetCodeOutput) => {
            this.codeString = resp.data
            this.loading = false
            this.initialLoad = false
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

    @Watch('dataId')
    refreshCode() {
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
        }

        if (this.dataId != -1) {
            params.dataId = this.dataId
        }

        if (this.scriptId != -1) {
            params.scriptId = this.scriptId
        }

        allCode(params).then((resp : TAllCodeOutput) => {
            this.allCode = resp.data
            if (this.allCode.length > 0) {
                this.selectedCode = this.allCode[0]
            } else {
                this.selectedCode = null

                // Force pull to get proper behavior as if
                // the selected code changed. This should
                // be cheap so even if the code DID change
                // not much should happen.
                this.pullCode()
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

    onSave() {
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
            code: this.codeString,
        }

        this.saveInProgress = true

        if (this.dataId != -1) {
            params.dataId = this.dataId
        }

        if (this.scriptId != -1) {
            params.scriptId = this.scriptId
        }

        saveCode(params).then((resp : TSaveCodeOutput) => {
            this.allCode.unshift(resp.data)
            this.selectedCode = resp.data
            this.saveInProgress = false
        }).catch((err : any) => {
            this.saveInProgress = false
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    onEvent(event : string, ...args : any[]) {
        this.$emit(event, ...args)
    }

    mounted() {
        this.refreshCode()
        //@ts-ignore
        let ele : HTMLElement = this.$refs.editor.$el
        // Add events here to let toolbar handle input events.
        document.addEventListener('keydown', (e : KeyboardEvent) => {
            if (!document.activeElement) {
                return
            }
            
            // This needs to be here so that the delete doesn't
            // accidentally trigger a hotkey when a dialog is
            // in focus.
            if (!ele.contains(document.activeElement)) {
                return
            }

            this.$refs.toolbar.handleHotkeys(e)
        })
    }
}

</script>

<style scoped>

>>>.v-select__selections input {
    width: 30px;
}

</style>
