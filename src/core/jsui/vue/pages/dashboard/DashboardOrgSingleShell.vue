<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <v-overlay :value="!isReady">
                <v-progress-circular indeterminate size="64"></v-progress-circular>
            </v-overlay>

            <div v-if="isReady">
                <v-list-item class="px-4 pt-4">
                    <v-list-item-content>
                        <v-list-item-title class="title">
                            {{shellTypeStr }} Script: {{ script.Name }}
                        </v-list-item-title>
                    </v-list-item-content>

                    <v-dialog v-model="showHideDelete"
                              persistent
                              max-width="40%"
                    >
                        <template v-slot:activator="{ on }">
                            <v-btn color="error" v-on="on">
                                Delete
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            item-name="shell scripts"
                            :items-to-delete="[script.Name]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item>
                <v-divider></v-divider>

                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row class="mx-4">
                            <v-col cols="6">
                                <create-new-shell-script-form
                                    edit-mode
                                    :shell-type="script.TypeId"
                                    :reference-script="script"
                                    @do-save="onEdit"
                                >
                                </create-new-shell-script-form>
                            </v-col>

                            <v-col cols="6">
                                <v-card>
                                    <v-list-item>
                                        <v-list-item-content
                                            class="title"
                                            style="flex: 1 1 !important;"
                                        >
                                            Versions
                                        </v-list-item-content>

                                        <v-list-item-content
                                            style="flex: 4 1 !important;"
                                        >
                                            <v-select
                                                v-model="selectedVersion"
                                                :items="versionItems"
                                                class="ml-4"
                                                filled
                                                hide-details
                                                dense
                                            >
                                            </v-select>
                                        </v-list-item-content>

                                        <v-list-item-action>
                                            <v-dialog
                                                v-model="showHideRun"
                                                persistent
                                                max-width="40%"
                                            >
                                                <template v-slot:activator="{on}">
                                                    <v-btn
                                                        color="primary"
                                                        icon
                                                        x-small
                                                        v-on="on"
                                                    >
                                                        <v-icon>mdi-play</v-icon>
                                                    </v-btn>
                                                </template>

                                                <v-card>
                                                    <server-table-with-controls
                                                        class="ma-4"
                                                        disable-new
                                                        disable-delete
                                                        enable-select
                                                        v-model="serversToRun"
                                                    >
                                                    </server-table-with-controls>

                                                    <v-card-actions>
                                                        <v-btn
                                                            color="error"
                                                            @click="showHideRun=false"
                                                        >
                                                            Cancel
                                                        </v-btn>

                                                        <v-spacer></v-spacer>

                                                        <v-btn
                                                            color="success"
                                                            @click="requestRunOnServers"
                                                            :loading="requestRunInProgress"
                                                        >
                                                            Run
                                                        </v-btn>
                                                    </v-card-actions>
                                                </v-card>
                                            </v-dialog>
                                        </v-list-item-action>
                                    </v-list-item>

                                    <v-divider></v-divider>

                                    <generic-code-editor
                                        v-if="!!currentScriptText"
                                        v-model="currentScriptText"
                                        :lang="shellCMLanguage"
                                        :readonly="!canEditScript"
                                        full-height
                                        :height-offset="-52"
                                    >
                                    </generic-code-editor>

                                    <v-row v-else justify="center" align="center">
                                        <v-progress-circular indeterminate size="64"></v-progress-circular>
                                    </v-row>

                                    <v-card-actions v-if="!!currentScriptText">
                                        <v-btn
                                            color="error"
                                            @click="cancelEditScript"
                                            v-if="canEditScript"
                                        >
                                            Cancel
                                        </v-btn>
                                        <v-spacer></v-spacer>
                                        <v-btn
                                            color="success"
                                            @click="saveScript"
                                            v-if="canEditScript"
                                            :loading="saveInProgress"
                                        >
                                            Save
                                        </v-btn>

                                        <v-btn
                                            color="success"
                                            @click="canEditScript = true"
                                            v-if="!canEditScript"
                                            :disabled="!currentScriptText"
                                            :loading="loadInProgress"
                                        >
                                            Edit
                                        </v-btn>
                                    </v-card-actions>
                                </v-card>
                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Runs</v-tab>
                    <v-tab-item>
                        <shell-run-table-with-controls
                            :shell-id="script.Id"
                        >
                        </shell-run-table-with-controls>
                    </v-tab-item>
                </v-tabs>
            </div>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import CreateNewShellScriptForm from '../../components/dashboard/CreateNewShellScriptForm.vue'
import GenericCodeEditor from '../../generic/code/GenericCodeEditor.vue'
import GenericDeleteConfirmationForm from '../../components/dashboard/GenericDeleteConfirmationForm.vue'
import MetadataStore from '../../../ts/metadata'
import ServerTableWithControls from '../../generic/resources/ServerTableWithControls.vue'
import ShellRunTableWithControls from '../../generic/resources/ShellRunTableWithControls.vue'
import { createUserString } from '../../../ts/users'
import { PageParamsStore } from '../../../ts/pageParams'
import {
    getShellScript, TGetShellScriptOutput,
    allShellScriptVersions, TAllShellScriptVersionsOutput,
    getShellScriptVersion, TGetShellScriptVersionOutput,
    newShellScriptVersion, TNewShellScriptVersionOutput,
    deleteShellScript,
} from '../../../ts/api/apiShell'
import {
    requestRunShellScript, TRequestRunShellScriptOutput
} from '../../../ts/api/apiShellRun'
import { contactUsUrl, createOrgShellUrl, createSingleShellRequestUrl } from '../../../ts/url'
import {
    ShellTypes,
    ShellScript,
    ShellScriptVersion,
    ShellTypeToCodeMirror,
} from '../../../ts/shell'
import { standardFormatTime } from '../../../ts/time'
import { Server } from '../../../ts/infrastructure'

@Component({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        CreateNewShellScriptForm,
        GenericCodeEditor,
        GenericDeleteConfirmationForm,
        ServerTableWithControls,
        ShellRunTableWithControls
    },
})

export default class DashboardOrgSingleShell extends Vue {
    script : ShellScript | null = null
    versions : ShellScriptVersion[] | null = null

    selectedVersion : ShellScriptVersion | null = null

    currentScriptText : string | null = null
    backupScriptText:  string = ""

    canEditScript: boolean = false
    showHideDelete : boolean = false
    showHideRun: boolean = false

    saveInProgress : boolean = false
    loadInProgress : boolean = false

    serversToRun : Server[] = []
    requestRunInProgress : boolean = false

    get isReady() : boolean {
        return !!this.script && !!this.versions
    }

    get shellTypeStr() : string {
        return ShellTypes[this.script!.TypeId]
    }

    get shellCMLanguage() : string {
        return ShellTypeToCodeMirror.get(<ShellTypes>this.script!.TypeId)!
    }

    get versionItems() : any[] {
        return this.versions!.map((ele : ShellScriptVersion, idx : number) => {
            let user = createUserString(MetadataStore.getters.getUser(ele.UploadUserId))
            return {
                text: `#${this.versions!.length - idx} ${user} [${standardFormatTime(ele.UploadTime)}]`,
                value: ele,
            }
        })
    }

    @Watch('selectedVersion')
    onChangeVersion() {
        this.canEditScript = false

        if (!this.selectedVersion) {
            return
        }

        this.loadInProgress = true
        getShellScriptVersion({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: this.script!.Id,
            version: this.selectedVersion!.Id,
        }).then((resp : TGetShellScriptVersionOutput) => {
            this.currentScriptText = resp.data
            this.backupScriptText = resp.data
            this.loadInProgress = false
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

    refreshData() {
        let id = parseInt(PageParamsStore.state.resource!.Id, 10)
        getShellScript({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: id,
        }).then((resp : TGetShellScriptOutput) => {
            this.script = resp.data

            // While this could theoretically be done in parallel,
            // we need this.script to be set before the version is set
            // as we need to know the script id for querying the version's
            // script text.
            allShellScriptVersions({
                orgId: PageParamsStore.state.organization!.Id,
                shellId: id,
            }).then((resp : TAllShellScriptVersionsOutput) => {
                this.versions = resp.data
                if (this.versions!.length > 0) {
                    
                    let urlParams = new URLSearchParams(window.location.search)
                    if (urlParams.has("version")) {
                        this.selectedVersion = this.versions![
                            Math.min(
                                Math.max(
                                    this.versions!.length - Number(urlParams.get("version")),
                                    0),
                                this.versions!.length - 1,
                            )
                        ]
                    } else {
                        this.selectedVersion = this.versions![0]
                    }
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
        this.refreshData()
    }

    onEdit(s : ShellScript) {
        this.script = s
    }

    cancelEditScript() {
        this.currentScriptText = this.backupScriptText
        this.canEditScript = false
    }

    saveScript() {
        this.saveInProgress = true
        newShellScriptVersion({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: this.script!.Id,
            script: this.currentScriptText!,
        }).then((resp : TNewShellScriptVersionOutput) => {
            this.versions!.unshift(resp.data)
            this.selectedVersion = resp.data
            this.canEditScript = false
            this.saveInProgress = false
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

    onDelete() {
        deleteShellScript({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: this.script!.Id,
        }).then(() => {
            window.location.assign(createOrgShellUrl(
                PageParamsStore.state.organization!.OktaGroupName
            ))
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

    requestRunOnServers() {
        this.requestRunInProgress = true

        requestRunShellScript({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: this.script!.Id,
            versionId: this.selectedVersion!.Id,
            servers: this.serversToRun.map((ele : Server) => ele.Id),
        }).then((resp : TRequestRunShellScriptOutput) => {
            this.showHideRun = false
            this.serversToRun = []

            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Successfully submitted your run request.",
                true,
                "View",
                createSingleShellRequestUrl(
                    PageParamsStore.state.organization!.OktaGroupName,
                    resp.data.RequestId!,
                ),
                false);

        }).catch((err: any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.requestRunInProgress = false
        })
    }
}

</script>
