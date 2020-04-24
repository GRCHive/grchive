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
                        Script Run Request:
                        </span>
                        <span class="title font-weight-regular">
                        {{ req.Name }}
                        </span>
                    </v-list-item-title>

                    <v-list-item-subtitle style="display: flex; align-items: center;">
                        <span>Script: </span>
                        <a :href="scriptUrl">{{ script.Name }}</a>
                        <hash-renderer
                            class="ml-4"
                            :hash="code.GitHash"
                        >
                        </hash-renderer>
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-content v-if="!!oneTime || !schedule" style="flex-grow: 1;">
                    <span class="font-weight-bold">When: </span>

                    <span v-if="!!oneTime">
                        {{ oneTimeStr }}
                    </span>

                    <span v-else>
                        Immediately Upon Approval
                    </span>
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
                            item-name="script run requests"
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

                            <v-card class="mt-4" v-if="!!schedule">
                                <v-card-title>
                                    Script Schedule
                                </v-card-title>
                                <v-divider></v-divider>

                                <create-scheduled-event-form
                                    class="ma-4"
                                    :value="schedule"
                                    no-name
                                    readonly
                                >
                                </create-scheduled-event-form>

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

                <v-tab>Script</v-tab>
                <v-tab-item>
                    <managed-code-ide
                        :script-id="script.Id"
                        :code-id="code.Id"
                        :initial-params="runParams"
                        readonly
                        disable-run
                        disable-save
                        lang="text/x-kotlin"
                        full-height
                    >
                    </managed-code-ide>
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
    getGenericRequestScript, TGetGenericRequestScriptOutput,
    getGenericRequest, TGetGenericRequestOutput,
    editGenericRequest,
    approveDenyGenericRequest, TApproveDenyRequestOutput,
    deleteGenericRequest
} from '../../../ts/api/apiRequests'
import {
    contactUsUrl,
    createSingleScriptUrl,
    createOrgDocRequestsUrl,
} from '../../../ts/url'
import { ClientScript } from '../../../ts/clientScripts'
import { ManagedCode } from '../../../ts/code'
import { GenericRequest, GenericApproval } from '../../../ts/requests'
import HashRenderer from '../../generic/code/HashRenderer.vue'
import GenericDeleteConfirmationForm from '../../components/dashboard/GenericDeleteConfirmationForm.vue'
import CreateNewGenericRequestForm from '../../components/dashboard/CreateNewGenericRequestForm.vue'
import GenericApprovalDisplay from '../../generic/requests/GenericApprovalDisplay.vue'
import ManagedCodeIde from '../../generic/code/ManagedCodeIDE.vue'
import { standardFormatTime } from '../../../ts/time'
import {
    ScheduledEvent,
    createScheduledEventFromRRule
} from '../../../ts/event'
import CreateScheduledEventForm from '../../generic/CreateScheduledEventForm.vue'
import CommentManager from '../../generic/CommentManager.vue'

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
        GenericDeleteConfirmationForm,
        HashRenderer,
        CreateNewGenericRequestForm,
        GenericApprovalDisplay,
        ManagedCodeIde,
        CreateScheduledEventForm,
        CommentManager,
    }
})
export default class DashboardOrgSingleScriptRequest extends Vue {
    req : GenericRequest | null | undefined = null
    approval : GenericApproval | null | undefined = null
    script : ClientScript | null = null
    code : ManagedCode | null = null
    oneTime : Date | null | undefined = null
    schedule : ScheduledEvent | null = null

    runParams : Record<string, any> = Object()

    showHideDelete : boolean = false

    canEditRequest : boolean = false
    requestValid : boolean = false
    editInProgress: boolean = false

    showHideDenyReason : boolean = false
    denyReason : string = ""

    $refs!: {
        editForm: CreateNewGenericRequestForm
    }

    get oneTimeStr() : string {
        return standardFormatTime(this.oneTime!)
    }

    get isLoading() : boolean {
        return (!this.script || !this.code || !this.req)
    }

    get scriptUrl() : string {
        return createSingleScriptUrl(PageParamsStore.state.organization!.OktaGroupName, this.script!.Id)
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

        getGenericRequestScript({
            orgId: PageParamsStore.state.organization!.Id,
            requestId: Number(PageParamsStore.state.resource!.Id),
        }).then((resp : TGetGenericRequestScriptOutput) => {
            this.script = resp.data.Script
            this.code = resp.data.Code
            this.oneTime = resp.data.OneTime
            this.schedule = createScheduledEventFromRRule(resp.data.RRule)
            this.runParams = resp.data.Params
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
        approveDenyGenericRequest({
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
