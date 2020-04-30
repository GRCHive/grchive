<template>
    <div>
        <v-dialog v-model="showHideConfirmRunSave" persistent max-width="40%">
            <v-card>
                <v-card-title>
                    Did you mean to save?
                </v-card-title>
                <v-divider></v-divider>

                <div class="ma-2">
                    Running this script at latest without saving will not pick up your latest changes. If that is not what you meant to do, please hit Cancel and Save before running the script again.
                </div>

                <v-divider></v-divider>
                <v-card-actions>
                    <v-btn color="primary" @click="showHideConfirmRunSave = false">
                        Cancel
                    </v-btn>
                    <v-spacer></v-spacer>
                    <v-btn color="error" @click="run(selectedCode, true)">
                        Continue
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>

        <generic-code-toolbar
            @save="onSave"
            @revert="pullCode"
            ref="toolbar"
            :save-in-progress="saveInProgress"
            :disable-save="disableSave"
        >
            <template v-slot:custom-menu>
                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn text color="accent" v-on="on">
                            View
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
                <v-icon color="warning" v-if="dirty" id="saveIcon">
                    mdi-content-save
                </v-icon>

                <v-col cols="auto">
                    <v-select
                        dense
                        solo
                        flat
                        hide-details
                        :items="allCodeItems"
                        label="Version"
                        v-model="selectedCode"
                        :readonly="readonly"
                    >
                        <template v-slot:item="{ item }">
                            <hash-renderer
                                :hash="item.value.GitHash"
                                class="mr-2"
                            >
                            </hash-renderer>
                            <code-build-status
                                :commit="item.value.GitHash"
                                class="mr-2"
                            >
                            </code-build-status>
                            {{ item.text }} 
                        </template>

                        <template v-slot:selection="{ item }">
                            <hash-renderer
                                :hash="item.value.GitHash"
                                class="mr-2"
                            >
                            </hash-renderer>
                            <code-build-status
                                :commit="item.value.GitHash"
                                class="mr-2"
                            >
                            </code-build-status>
                            {{ item.text }} 
                        </template>
                    </v-select>
                </v-col>

                <v-icon>
                    {{ finalReadonly ? 'mdi-eye' : 'mdi-pencil' }}
                </v-icon>
            </template>
        </generic-code-toolbar>
        <v-divider></v-divider>

        <v-row>
            <v-col :cols="scriptId != -1 ? 9 : 12">
                <dynamic-split-container :enable-col-b="showLogs">
                    <template v-slot:first-col>
                        <generic-code-editor
                            :value="codeString"
                            :lang="lang"
                            :readonly="finalReadonly"
                            :full-height="fullHeight"
                            @input="onInput"
                            ref="editor"
                        >
                        </generic-code-editor>
                    </template>

                    <template v-slot:second-col>
                        <log-viewer
                            :commit="selectedCode.GitHash"
                            full-height
                        >
                        </log-viewer>
                    </template>
                </dynamic-split-container>
            </v-col>

            <v-col cols="3" v-if="scriptId != -1" class="pr-2 py-0 pl-0">
                <script-params-editor
                    :linked-client-data.sync="currentLinkedClientData"
                    :script-parameter-types.sync="currentScriptParams"
                    :script-parameter-values.sync="currentScriptParamValues" 
                    :readonly="finalReadonly"
                    :run-in-progress="runInProgress"
                    @runLatest="runScriptLatest"
                    @runRevision="runScriptAtRevision"
                    @scheduleRun="scheduleScript"
                    :disableRun="!selectedCode || disableRun"
                >
                </script-params-editor>
            </v-col>
        </v-row>
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
    runCode, TRunCodeOutput,
} from '../../../ts/api/apiCode'
import {
    contactUsUrl,
    createSingleScriptRequestUrl
} from '../../../ts/url'
import { standardFormatTime } from '../../../ts/time'
import { FullClientDataWithLink } from '../../../ts/clientData'
import {
    CodeParamType
} from '../../../ts/code'
import { ScheduledEvent } from '../../../ts/event'

import LogViewer from '../logs/LogViewer.vue'
import DynamicSplitContainer from '../DynamicSplitContainer.vue'
import CodeBuildStatus from './CodeBuildStatus.vue'
import ScriptParamsEditor from './ScriptParamsEditor.vue'
import HashRenderer from './HashRenderer.vue'

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
        codeId: {
            type: Number,
            default: -1,
        },
        disableRun: {
            type: Boolean,
            default: false,
        },
        disableSave: {
            type: Boolean,
            default: false,
        },
        initialParams: {
            type: Object,
            default: () => Object() as Record<string, any>
        }
    }
})

@Component({
    components: {
        LogViewer,
        DynamicSplitContainer,
        GenericCodeToolbar,
        GenericCodeEditor,
        CodeBuildStatus,
        ScriptParamsEditor,
        HashRenderer,
    }
})
export default class ManagedCodeIDE extends mixins(Props, ManagedProps) {
    codeString : string = ""

    dirty : boolean = false

    loading: boolean = false
    showLogs: boolean = false

    showHideConfirmRunSave : boolean = false

    // This is equivalent to loading for the first time
    // we load code. We need this because the code toolbar
    // will use 'codeString' when mounted to determine if the user
    // needs to save. Thus, we need to delay its mount until after
    // the code loads. We don't want to use 'loading' to prevent
    // the mount because otherwise there's a flicker every time you
    // change the code version which is slightly unpleasant.
    initialLoad: boolean = true

    saveInProgress: boolean = false
    runInProgress: boolean = false

    allCode : ManagedCode[] = []
    selectedCode : ManagedCode | null = null

    currentLinkedClientData : FullClientDataWithLink[] = []
    currentScriptParams : (CodeParamType | null)[] = []
    currentScriptParamValues : Record<string, any> = Object()

    $refs!: {
        editor: GenericCodeEditor,
        toolbar: GenericCodeToolbar,
    }

    @Watch('codeString')
    @Watch('currentLinkedClientData')
    @Watch('currentScriptParams')
    markDirty() {
        this.dirty = true
    }

    onInput(text : string) {
        this.codeString = text
    }

    get finalReadonly() : boolean {
        // Make only the latest revision non readonly.
        // This way the "Run at Revision" functionality becomes more clear as they won't
        // be able to have changes to it.
        return this.readonly || (!this.selectedCode && this.allCode.length != 0) ||
            (this.allCode.length > 0 && this.allCode[0].Id != this.selectedCode!.Id)
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
            this.codeString = resp.data.Code

            if (!!resp.data.ScriptData) {
                this.currentLinkedClientData = resp.data.ScriptData.ClientData
                this.currentScriptParams = resp.data.ScriptData.Params
            }

            this.loading = false
            this.initialLoad = false

            // A hack to make sure the initial code/params pull doesn't make
            // us as being dirty. Use two nextTick() here because I think the
            // @Watch also triggers on the next tick so resetting dirty needs
            // to take place after two ticks.
            Vue.nextTick(() => {
                Vue.nextTick(() => {
                    this.dirty = false
                })
            })
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
                if (this.codeId == -1) {
                    this.selectedCode = this.allCode[0]
                } else {
                    this.selectedCode = this.allCode.find((ele : ManagedCode) => ele.Id == this.codeId)!
                }
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
        if (this.disableSave) {
            return
        }

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
            params.scriptData = {
                params: this.currentScriptParams,
                clientDataId: this.currentLinkedClientData.map((ele : FullClientDataWithLink) => ele.Data.Id),
            }
        }

        saveCode(params).then((resp : TSaveCodeOutput) => {
            this.allCode.unshift(resp.data)
            this.selectedCode = resp.data
            this.saveInProgress = false
            this.dirty = false
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

    handleUnload(e : Event) {
        if (this.dirty) {
            e.preventDefault()
            e.returnValue = false
        }
    }

    mounted() {
        this.currentScriptParamValues = this.initialParams
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

        window.addEventListener('beforeunload', this.handleUnload)
    }

    runScriptLatest() {
        if (this.disableRun) {
            return
        }

        if (this.dirty) {
            // Need to make sure the user knows that their latest changes won't be picked up
            // until they save.
            this.showHideConfirmRunSave = true
            return
        }

        this.run(this.selectedCode, true)
    }

    runScriptAtRevision() {
        if (this.disableRun) {
            return
        }
        this.run(this.selectedCode, false)
    }

    run(code : ManagedCode | null, latest : boolean, schedule : ScheduledEvent | null = null) {
        if (this.disableRun) {
            return
        }

        if (!code) {
            return
        }

        this.runInProgress = true

        runCode({
            orgId: PageParamsStore.state.organization!.Id,
            codeId: code!.Id,
            latest: latest,
            params: this.currentScriptParamValues,
            schedule: schedule,
        }).then((resp : TRunCodeOutput) => {
            this.runInProgress = false
            if (!!schedule) {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Your run request has been successfully submitted.",
                    false,
                    "Track",
                    createSingleScriptRequestUrl(
                        PageParamsStore.state.organization!.OktaGroupName,
                        resp.data,
                    ),
                    false);
            } else {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Your run request has been successfully submitted.",
                    true,
                    "Track",
                    createSingleScriptRequestUrl(
                        PageParamsStore.state.organization!.OktaGroupName,
                        resp.data,
                    ),
                    false);
            }
        }).catch((err : any) => {
            this.runInProgress = false
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    scheduleScript(s : ScheduledEvent) {
        if (this.disableRun) {
            return
        }
        this.run(this.selectedCode, false, s)
    }
}

</script>

<style scoped>

>>>.v-select__selections input {
    width: 30px;
}

#saveIcon {
    animation-duration: 1s;
    animation-name: blinkSave;
    animation-iteration-count: infinite;
    animation-direction: alternate;
}

@keyframes blinkSave {
    from {
        opacity: 1;
    }

    to {
        opacity: 0.1;
    }
}

</style>
