<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        PBC Request: {{ currentRequest.Name }}
                        <doc-request-status-display
                            :status="status"
                        >
                        </doc-request-status-display>
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        <p class="ma-0">
                            <span class="font-weight-bold">Relevant Document Category:</span>
                            <a :href="parentCategoryUrl">{{ parentCategory.Name }}</a>
                        </p>

                        <p class="ma-0" v-if="!!parentControl && !!controlFolder">
                            <span class="font-weight-bold">Relevant Control:</span>
                            <a :href="parentControlUrl">{{ parentControl.Name }}</a>

                            <span class="font-weight-bold">, Folder:</span>
                            <span>{{ controlFolder.Name }}</span>
                        </p>

                        <p class="ma-0" v-if="isApproved">
                            <span class="font-weight-bold">Approved By:</span>
                            <span>{{ approveUserName }}</span>
                            <span> at </span>
                            <span>{{ approveTime }}</span>
                        </p>
                    </v-list-item-subtitle>

                </v-list-item-content>

                <v-list-item-action>
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
                            item-name="document requests"
                            :items-to-delete="[currentRequest.Name]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>

                </v-list-item-action>

                <v-list-item-action v-if="canComplete">
                    <v-btn 
                        @click="doComplete(true)"
                        color="success"
                    >
                        Complete  
                    </v-btn>
                </v-list-item-action>

                <v-list-item-action v-if="canReopen">
                    <v-btn 
                        @click="doComplete(false)"
                        color="warning"
                        :class="isApproved ? '' : 'ml-4'"
                    >
                        Reopen
                    </v-btn>
                </v-list-item-action>

                <v-list-item-action v-if="canApprove">
                    <v-btn 
                        @click="doApprove"
                        color="success"
                    >
                        Approve
                    </v-btn>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="6">
                                <create-new-request-form
                                    class="mb-4"
                                    ref="form"
                                    edit-mode
                                    load-cats
                                    :reference-cat="parentCategory"
                                    :reference-req="currentRequest"
                                    :reference-control="parentControl"
                                    :reference-folder="controlFolder"
                                    :vendor-product-id="vendorProductId"
                                    @do-save="onEdit"
                                ></create-new-request-form>

                                <v-card class="mb-4">
                                    <v-card-title>
                                        Requester Information
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <div class="ma-4">
                                        <user-search-form-component
                                            label="Requester"
                                            :user="requestUser"
                                            readonly
                                        >
                                        </user-search-form-component>

                                        <v-text-field
                                            :value="requestTime"
                                            label="Request Time"
                                            prepend-icon="mdi-calendar"
                                            readonly
                                        >
                                        </v-text-field>
                                    </div>
                                </v-card>

                                <v-card>
                                    <v-card-title>
                                        Relevant Files
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <doc-file-manager
                                        :cat-id="!!parentCategory ? parentCategory.Id : -1"
                                        :request-id="currentRequest.Id"
                                        :request-linked-to-control="!!parentControl"
                                        :request-control="parentControl"
                                        :vendor-id="vendorId"
                                        :vendor-product-id="vendorProductId"
                                        :folder="controlFolder"
                                        disable-sample
                                        disable-delete
                                        @changed-all-files="onChangedFiles"
                                    ></doc-file-manager>
                                </v-card>
                            </v-col>

                            <v-col cols="6">
                                <v-card>
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

                    <v-tab>Audit Trail</v-tab>
                    <v-tab-item>
                        <audit-trail-viewer
                            :resource-type="['document_requests']"
                            :resource-id="[`${currentRequest.Id}`]"
                            no-header
                        >
                        </audit-trail-viewer>
                    </v-tab-item>

                </v-tabs>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { DocumentRequest, DocRequestStatus, getDocumentRequestStatus } from '../../../ts/docRequests'
import { createUserString } from '../../../ts/users'
import { TGetSingleDocumentRequestOutput, getSingleDocRequest } from '../../../ts/api/apiDocRequests'
import { 
    deleteSingleDocRequest,
    completeDocRequest,
    reopenDocRequest,
    approveDocRequest,
} from '../../../ts/api/apiDocRequests'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgDocRequestsUrl, createSingleDocCatUrl, createControlUrl } from '../../../ts/url'
import { ControlDocumentationCategory, ControlDocumentationFile } from '../../../ts/controls'
import MetadataStore from '../../../ts/metadata'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CreateNewRequestForm from './CreateNewRequestForm.vue'
import DocFileManager from '../../generic/DocFileManager.vue'
import CommentManager from '../../generic/CommentManager.vue'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'
import { standardFormatTime } from '../../../ts/time'
import { allDocRequestDocCatLink, TAllDocRequestDocCatLinksOutput } from '../../../ts/api/apiDocRequestsDocCatLinks'
import { allDocRequestControlLink, TAllDocRequestControlLinksOutput } from '../../../ts/api/apiDocRequestControlLinks'
import { getDocRequestControlFolderLink, TGetDocRequestControlFolderLinksOutput } from '../../.././ts/api/apiDocRequestFolderLinks'
import { FileFolder } from '../../../ts/folders'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'
import DocRequestStatusDisplay from '../../generic/requests/DocRequestStatusDisplay.vue'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CreateNewRequestForm,
        DocFileManager,
        CommentManager,
        UserSearchFormComponent,
        AuditTrailViewer,
        DocRequestStatusDisplay,
    }
})
export default class FullEditDocRequestComponent extends Vue {
    currentRequest : DocumentRequest | null = null
    relevantFiles: ControlDocumentationFile[] = []
    vendorId: number = -1
    vendorProductId: number = -1

    parentCategory : ControlDocumentationCategory | null = null
    parentControl : ProcessFlowControl | null = null
    controlFolder : FileFolder | null = null

    showHideDelete: boolean = false

    get status() : DocRequestStatus {
        return getDocumentRequestStatus(this.currentRequest!)
    }

    get canComplete() : boolean {
        return this.status == DocRequestStatus.Open ||
            this.status == DocRequestStatus.InProgress ||
            this.status == DocRequestStatus.Feedback ||
            this.status == DocRequestStatus.Overdue
    }

    get canReopen() : boolean {
        return !this.canComplete
    }

    get canApprove() : boolean {
        return !this.canComplete && this.status != DocRequestStatus.Approved
    }

    get commentParams() : Object {
        return {
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }
    }

    get ready() : boolean {
        return !!this.currentRequest && !!this.parentCategory
    }

    get isApproved() : boolean { 
        return !!this.currentRequest && !!this.currentRequest.ApproveTime && !!this.currentRequest.ApproveUserId
    }

    get parentControlUrl() : string {
        if (!this.parentControl) {
            return ""
        }
        return createControlUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            this.parentControl!.Id)
    }

    get parentCategoryUrl() : string {
        if (!this.parentCategory) {
            return ""
        }
        return createSingleDocCatUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            this.parentCategory!.Id)
    }

    get requestUser() : User | null {
        if (!this.currentRequest) {
            return null
        }
        return MetadataStore.getters.getUser(this.currentRequest.RequestedUserId)
    }

    get requestTime() : string {
        if (!this.currentRequest) {
            return ""
        }
        return standardFormatTime(this.currentRequest.RequestTime)
    }

    get completionTime() : string {
        if (!this.currentRequest || !this.currentRequest.CompletionTime) {
            return ""
        }
        return standardFormatTime(this.currentRequest.CompletionTime)
    }

    get approveUserName() : string {
        if (!this.currentRequest) {
            return ""
        }
        let user = MetadataStore.getters.getUser(this.currentRequest.ApproveUserId)
        return createUserString(user)
    }

    get approveTime() : string {
        if (!this.currentRequest || !this.currentRequest.ApproveTime) {
            return ""
        }
        return standardFormatTime(this.currentRequest.ApproveTime)
    }

    onError(err : any) { 
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops! Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    refreshParentControl() {
        allDocRequestControlLink({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllDocRequestControlLinksOutput) => {
            this.parentControl = resp.data.Control!

            if (!!this.parentControl) {
                getDocRequestControlFolderLink({
                    requestId: this.currentRequest!.Id,
                    orgId: PageParamsStore.state.organization!.Id,
                    controlId: this.parentControl.Id,
                }).then((resp : TGetDocRequestControlFolderLinksOutput) => {
                    this.controlFolder = resp.data
                }).catch(this.onError)
            }
        }).catch(this.onError)
    }

    refreshParentCategory() {
        allDocRequestDocCatLink({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllDocRequestDocCatLinksOutput) => {
            this.parentCategory = resp.data.Cat!
        }).catch(this.onError)
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getSingleDocRequest({
            requestId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSingleDocumentRequestOutput) => {
            this.currentRequest = resp.data.Request
            this.relevantFiles = resp.data.Files
            this.vendorProductId = resp.data.VendorProductId
            this.vendorId = resp.data.VendorId

            this.refreshParentCategory()
            this.refreshParentControl()
        }).catch(this.onError)
    }

    onEdit(req : DocumentRequest) {
        this.currentRequest = req
    }

    onDelete() {
        deleteSingleDocRequest({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
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

    doComplete(val : boolean) {
        if (val) {
            completeDocRequest({
                requestId: this.currentRequest!.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then(() => {
                this.currentRequest!.CompletionTime = new Date()
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        } else {
            reopenDocRequest({
                requestId: this.currentRequest!.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then(() => {
                this.currentRequest!.FeedbackTime = new Date()
                this.currentRequest!.ApproveTime = null
                this.currentRequest!.ApproveUserId = null
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
    }

    doApprove() {
        approveDocRequest({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.currentRequest!.ApproveTime = new Date()
            this.currentRequest!.ApproveUserId = PageParamsStore.state.user!.Id
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

    onChangedFiles() {
        this.currentRequest!.ProgressTime = new Date()
    }

    mounted() {
        this.refreshData()
    }
}

</script>
