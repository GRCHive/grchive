<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-overlay :value="isLoading">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <v-content v-if="!isLoading">
            <v-list-item two-line class="px-4">
                <v-list-item-content style="flex-grow: 3;">
                    <v-list-item-title>
                        <span class="title">
                        Shell Script Run Request:
                        </span>
                        <span class="title font-weight-regular">
                        {{ req.Name }}
                        </span>
                    </v-list-item-title>

                    <v-list-item-subtitle style="display: flex; align-items: center;">
                        <span>{{ shellTypeStr }} Script: </span>
                        <a :href="shellUrl">{{ shell.Name }}</a>
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-spacer></v-spacer>

                <v-list-item-action v-if="!approval">
                    <v-dialog v-model="showHideDelete"
                              persistent
                              max-width="40%"
                    >
                        <template v-slot:activator="{ on }">
                            <v-btn color="warning" v-on="on">
                                Delete
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            item-name="shell script run requests"
                            :items-to-delete="[req.Name]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

                <v-list-item-action class="ml-4" v-if="!approval">
                    <v-dialog persistent max-width="40%" v-model="showHideDenyReason">
                        <template v-slot:activator="{on}">
                            <v-btn 
                                color="error"
                                v-on="on"
                            >
                                Deny
                            </v-btn>
                        </template>

                        <v-card>
                            <v-card-title>
                                Denial Reason
                            </v-card-title>
                            <v-divider></v-divider>

                            <div class="ma-4">
                                <v-textarea v-model="denyReason"
                                            label="Reason"
                                            filled
                                            hide-details
                                ></v-textarea> 
                            </div>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="showHideDenyReason = false"
                                >
                                    Cancel
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-btn
                                    color="success"
                                    @click="onApproveDeny(false, denyReason)"
                                >
                                    Save
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </v-list-item-action>

                <v-list-item-action v-if="!approval">
                    <v-btn
                        color="success"
                        @click="onApproveDeny(true)"
                    >
                        Approve
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
            <v-divider></v-divider>

            <v-tabs>
                <v-tab>Overview</v-tab>
                <v-tab-item>
                    <v-row class="mx-4">
                        <v-col cols="7">
                            <v-card>
                                <v-card-title>
                                    Request
                                </v-card-title>
                                <v-divider></v-divider>

                                <create-new-generic-request-form
                                    ref="editForm"
                                    class="ma-4"
                                    v-model="req"
                                    :valid.sync="requestValid"
                                    :readonly="!canEditRequest"
                                >
                                </create-new-generic-request-form>

                                <v-card-actions>
                                    <v-btn
                                        color="error"
                                        @click="cancelEditRequest"
                                        v-if="canEditRequest"
                                    >
                                        Cancel
                                    </v-btn>
                                    <v-spacer></v-spacer>
                                    <v-btn
                                        color="success"
                                        @click="saveEditRequest"
                                        :disabled="!requestValid"
                                        v-if="canEditRequest"
                                    >
                                        Save
                                    </v-btn>

                                    <v-btn
                                        color="success"
                                        @click="canEditRequest = true"
                                        v-if="!canEditRequest"
                                    >
                                        Edit
                                    </v-btn>
                                </v-card-actions>
                            </v-card>
                        </v-col>

                        <v-col cols="5">
                            <generic-approval-display
                                :approval="approval"
                            >
                            </generic-approval-display>

                            <v-card class="mt-4">
                                <v-card-title>
                                    Comments
                                </v-card-title>
                                <v-divider></v-divider>
                                <comment-manager
                                    :params="commentParams"
                                ></comment-manager>
                            </v-card>
                        </v-col>
                    </v-row>
                </v-tab-item>

                <v-tab>Shell Script</v-tab>
                <v-tab-item>
                    <generic-code-editor
                        v-if="!!currentScriptText"
                        v-model="currentScriptText"
                        :lang="shellCMLanguage"
                        readonly
                        full-height
                    >
                    </generic-code-editor>
                </v-tab-item>

                <v-tab>Servers</v-tab>
                <v-tab-item>
                    <server-table
                        :resources="relevantServers"
                    >
                    </server-table>
                </v-tab-item>
            </v-tabs>
        </v-content>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import {
    getGenericRequest, TGetGenericRequestOutput,
    editGenericRequest,
    deleteGenericRequest,
    getGenericRequestShell, TGetGenericRequestShellOutput,
    approveDenyShellRequest, TApproveDenyRequestOutput,
} from '../../../ts/api/apiRequests'
import {
    getShellScriptVersion, TGetShellScriptVersionOutput,
} from '../../../ts/api/apiShell'
import {
    contactUsUrl,
    createOrgDocRequestsUrl,
    createSingleShellUrl,
} from '../../../ts/url'
import {
    ShellScript,
    ShellScriptVersion,
    ShellTypeToCodeMirror,
    ShellTypes,
} from '../../../ts/shell'
import { Server } from '../../../ts/infrastructure'
import { GenericRequest, GenericApproval } from '../../../ts/requests'
import GenericDeleteConfirmationForm from '../../components/dashboard/GenericDeleteConfirmationForm.vue'
import CreateNewGenericRequestForm from '../../components/dashboard/CreateNewGenericRequestForm.vue'
import GenericApprovalDisplay from '../../generic/requests/GenericApprovalDisplay.vue'
import { standardFormatTime } from '../../../ts/time'
import {
    ScheduledEvent,
    createScheduledEventFromRRule
} from '../../../ts/event'
import CommentManager from '../../generic/CommentManager.vue'
import GenericCodeEditor from '../../generic/code/GenericCodeEditor.vue'
import ServerTable from '../../generic/ServerTable.vue'

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
        GenericDeleteConfirmationForm,
        CreateNewGenericRequestForm,
        GenericApprovalDisplay,
        CommentManager,
        GenericCodeEditor,
        ServerTable
    }
})
export default class DashboardOrgSingleShellRequest extends Vue {
    req : GenericRequest | null | undefined = null
    approval : GenericApproval | null | undefined = null

    showHideDelete : boolean = false

    canEditRequest : boolean = false
    requestValid : boolean = false
    editInProgress: boolean = false

    showHideDenyReason : boolean = false
    denyReason : string = ""

    shell : ShellScript | null = null
    shellVersion : ShellScriptVersion | null = null
    version : number = -1
    relevantServers: Server[] = []

    currentScriptText : string | null = null

    $refs!: {
        editForm: CreateNewGenericRequestForm
    }

    get shellTypeStr() : string {
        return ShellTypes[this.shell!.TypeId]
    }

    get shellCMLanguage() : string {
        return ShellTypeToCodeMirror.get(<ShellTypes>this.shell!.TypeId)!
    }

    get isLoading() : boolean {
        return (!this.req || !this.shell || !this.shellVersion)
    }

    get shellUrl() : string {
        return `${createSingleShellUrl(PageParamsStore.state.organization!.OktaGroupName, this.shell!.Id)}?version=${this.version}`
    }

    refreshData() {
        getGenericRequest({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: Number(PageParamsStore.state.resource!.Id),
        }).then((resp : TGetGenericRequestOutput) => {
            this.req = resp.data.Request
            this.approval = resp.data.Approval
        }).catch((err: any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })

        getGenericRequestShell({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: Number(PageParamsStore.state.resource!.Id),
        }).then((resp : TGetGenericRequestShellOutput) => {
            this.shell = resp.data.Shell
            this.shellVersion = resp.data.Version
            this.version = resp.data.VersionNum
            this.relevantServers = resp.data.Servers

            getShellScriptVersion({
                orgId: PageParamsStore.state.organization!.Id,
                shellId: this.shell!.Id,
                version: this.shellVersion!.Id,
            }).then((resp : TGetShellScriptVersionOutput) => {
                this.currentScriptText = resp.data
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

    onDelete() {
        deleteGenericRequest({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: this.req!.Id,
        }).then(() => {
            window.location.replace(createOrgDocRequestsUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onApproveDeny(approve : boolean, reason : string = "") {
        approveDenyShellRequest({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: this.req!.Id,
            approve: approve,
            reason: reason,
        }).then((resp : TApproveDenyRequestOutput) => {
            this.approval = resp.data
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

    cancelEditRequest() {
        this.canEditRequest = false
        this.$refs.editForm.resetToRefState()
    }

    saveEditRequest() {
        this.editInProgress = true
        editGenericRequest({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: this.req!.Id,
            request: this.req!,
        }).then(() => {
            this.editInProgress = false
            this.canEditRequest = false
            this.$refs.editForm.saveRefState()
        }).catch((err : any) => {
            this.editInProgress = false
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

    get commentParams() : Object {
        return {
            genericRequestId: this.req!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }
    }
}
</script>
