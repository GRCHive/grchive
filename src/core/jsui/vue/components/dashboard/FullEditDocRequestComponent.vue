<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Document Request: {{ currentRequest.Name }}
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        <p class="ma-0" v-if="!!parentCategory">
                            <span class="font-weight-bold">Relevant Document Category:</span>
                            <a :href="parentCategoryUrl">{{ parentCategory.Name }}</a>
                        </p>

                        <p class="ma-0" v-if="!!parentControl">
                            <span class="font-weight-bold">Relevant Control:</span>
                            <a :href="parentControlUrl">{{ parentControl.Name }}</a>
                        </p>

                        <p class="ma-0">
                            <span class="font-weight-bold">Status:</span>

                            <span v-if="!currentRequest.CompletionTime">
                                Pending
                                <v-icon
                                    small
                                    color="warning"
                                >
                                    mdi-help-circle
                                </v-icon>
                            </span>

                            <span v-else>
                                Complete ({{ completionTime }})
                                <v-icon
                                    small
                                    color="success"
                                >
                                    mdi-check
                                </v-icon>
                            </span>
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

                <v-list-item-action>
                    <v-btn 
                        @click="doComplete(true)"
                        color="success"
                        v-if="!currentRequest.CompletionTime"
                    >
                        Complete  
                    </v-btn>

                    <v-btn 
                        @click="doComplete(false)"
                        color="warning"
                        v-else
                    >
                        Reopen
                    </v-btn>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
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
                                v-model="relevantFiles"
                                disable-sample
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
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { DocumentRequest } from '../../../ts/docRequests'
import { TGetSingleDocumentRequestOutput, getSingleDocRequest } from '../../../ts/api/apiDocRequests'
import { deleteSingleDocRequest, completeDocRequest } from '../../../ts/api/apiDocRequests'
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

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CreateNewRequestForm,
        DocFileManager,
        CommentManager,
        UserSearchFormComponent
    }
})
export default class FullEditDocRequestComponent extends Vue {
    currentRequest : DocumentRequest | null = null
    relevantFiles: ControlDocumentationFile[] = []

    parentCategory : ControlDocumentationCategory | null = null
    parentControl : ProcessFlowControl | null = null

    showHideDelete: boolean = false

    $refs!: {
        form: CreateNewRequestForm
    }

    get commentParams() : Object {
        return {
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }
    }

    get ready() : boolean {
        return !!this.currentRequest
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

    @Watch('ready')
    onReady() {
        Vue.nextTick(() => {
            this.$refs.form.clearForm()
        })
    }

    refreshParentControl() {
        allDocRequestControlLink({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllDocRequestControlLinksOutput) => {
            this.parentControl = resp.data.Control!
        })
    }

    refreshParentCategory() {
        allDocRequestDocCatLink({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllDocRequestDocCatLinksOutput) => {
            this.parentCategory = resp.data.Cat!
        })
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

            this.refreshParentCategory()
            this.refreshParentControl()
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
        completeDocRequest({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            complete: val,
        }).then(() => {
            if (val) {
                this.currentRequest!.CompletionTime = new Date()
            } else {
                this.currentRequest!.CompletionTime = null
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
        this.refreshData()
    }
}

</script>
