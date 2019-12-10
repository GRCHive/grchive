<template>
    <div>
        <doc-file-table
            :resources="files"
            :selectable="true"
            :multi="true"
            v-model="selectedFiles"
        ></doc-file-table>

        <v-divider></v-divider>
        <v-list-item>
            <v-list-item-action class="ma-1">
                <v-btn icon @click="toggleSelection">
                    <v-icon v-if="!hasSelected">mdi-checkbox-blank-outline</v-icon>
                    <v-icon v-else>mdi-minus-box-outline</v-icon>
                </v-btn>
            </v-list-item-action>

            <v-dialog v-model="showHideDeleteFiles" persistent max-width="40%">
                <template v-slot:activator="{on}">
                    <v-btn color="error" v-on="on" :disabled="!hasSelected">
                        Delete
                    </v-btn>
                </template>

                <generic-delete-confirmation-form
                    item-name="documents"
                    :items-to-delete="selectedFileNames"
                    v-on:do-cancel="showHideDeleteFiles = false"
                    v-on:do-delete="deleteSelectedFiles"
                    :use-global-deletion="false"
                    :force-global-deletion="true">
                </generic-delete-confirmation-form>
            </v-dialog>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn color="info">
                    Sample
                </v-btn>
            </v-list-item-action>

            <v-list-item-action class="ml-4">
                <v-btn color="success" @click="downloadSelectedFiles" :disabled="!hasSelected">
                    Download
                </v-btn>
            </v-list-item-action>

            <v-list-item-action>
                <v-dialog v-model="showHideUpload" persistent max-width="40%">
                    <template v-slot:activator="{on}">
                        <v-btn color="primary" v-on="on">
                            Upload
                        </v-btn>
                    </template>
                    <upload-documentation-form :cat-id="catId"
                        @do-cancel="cancelUpload"
                        @do-save="finishUpload">
                    </upload-documentation-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

        <v-card>
            <v-card-title>
                <span class="mr-2">
                    Requests
                </span>
                <v-spacer></v-spacer>

                <v-dialog v-model="showHideRequest" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="warning" icon v-on="on">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <create-new-request-form
                        :cat-id="catId"
                        @do-cancel="showHideRequest = false"
                        @do-save="newRequest">
                    </create-new-request-form>
                </v-dialog>
            </v-card-title>
            <v-divider></v-divider>

            <doc-request-table :resources="requests">
            </doc-request-table>
        </v-card>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import UploadDocumentationForm from './UploadDocumentationForm.vue'
import { contactUsUrl } from '../../../ts/url'
import { ControlDocumentationFile } from '../../../ts/controls'
import { TGetControlDocumentsInput, TGetControlDocumentsOutput, getControlDocuments } from '../../../ts/api/apiControlDocumentation'
import { TDeleteControlDocumentsInput, TDeleteControlDocumentsOutput, deleteControlDocuments } from '../../../ts/api/apiControlDocumentation'
import { TDownloadControlDocumentsInput, TDownloadControlDocumentsOutput, downloadControlDocuments } from '../../../ts/api/apiControlDocumentation'
import { TGetAllDocumentRequestOutput, getAllDocRequests } from '../../../ts/api/apiDocRequests'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { saveAs } from 'file-saver'
import { PageParamsStore } from '../../../ts/pageParams'
import DocFileTable from '../../generic/DocFileTable.vue'
import CreateNewRequestForm from './CreateNewRequestForm.vue'
import { DocumentRequest } from '../../../ts/docRequests'
import DocRequestTable from '../../generic/DocRequestTable.vue'

export default Vue.extend({
    props : {
        catId: Number,
    },
    data : () => ({
        showHideUpload: false,
        pageNum: 0,
        totalPages: 0,
        files: [] as ControlDocumentationFile[],
        selectedFiles: [] as ControlDocumentationFile[],
        requests: [] as DocumentRequest[],
        showHideDeleteFiles: false,
        showHideRequest: false,
    }),
    components : {
        UploadDocumentationForm,
        GenericDeleteConfirmationForm,
        DocFileTable,
        CreateNewRequestForm,
        DocRequestTable
    },
    computed : {
        hasSelected() : boolean {
            return this.selectedFiles.length > 0
        },
        pageNumOneIndex() : number {
            if (this.totalPages == 0) {
                return 0
            }
            return this.pageNum + 1
        },
        selectedFileNames() : string[] {
            return this.selectedFiles.map((ele) => ele.StorageName)
        }
    },
    methods: {
        prepFile(f : ControlDocumentationFile) {
            f.RelevantTime = new Date(f.RelevantTime)
            f.UploadTime = new Date(f.UploadTime)
        },
        cancelUpload() {
            this.showHideUpload = false
        },
        finishUpload(newDoc : ControlDocumentationFile) {
            this.showHideUpload = false

            // We can avoid a refresh since we know it'll be the first item.
            this.prepFile(newDoc)
            this.files.unshift(newDoc)
        },
        refreshData(page : number, needPages : boolean) {
            getControlDocuments(<TGetControlDocumentsInput>{
                catId: this.catId,
                orgId: PageParamsStore.state.organization!.Id,
                page: page,
                needPages: needPages,
            }).then((resp : TGetControlDocumentsOutput) => {
                this.pageNum = resp.data.CurrentPage
                this.files = resp.data.Files
                for (let f of this.files) {
                    this.prepFile(f)
                }

                if (needPages) {
                    this.totalPages = resp.data.TotalPages
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

            getAllDocRequests({
                orgId: PageParamsStore.state.organization!.Id,
                catId: this.catId,
            }).then((resp : TGetAllDocumentRequestOutput) => {
                this.requests = resp.data
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        toggleSelection() {
            if (this.hasSelected) {
                this.selectedFiles = [] as ControlDocumentationFile[]
            } else {
                this.selectedFiles = this.files.slice()
            }
        },
        deleteSelectedFiles() {
            deleteControlDocuments(<TDeleteControlDocumentsInput>{
                orgId: PageParamsStore.state.organization!.Id,
                fileIds: this.selectedFiles.map((ele) => ele.Id),
                catId: this.catId,
            }).then(() => {
                let selectedFileSet = new Set(this.selectedFiles)
                for (let i = this.files.length - 1; i >= 0; --i) {
                    if (selectedFileSet.has(this.files[i])) {
                        this.files.splice(i, 1)
                    }
                }
                this.selectedFiles = []
                this.showHideDeleteFiles = false
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        downloadSelectedFiles() {
            // Download each file individually from the webserver and then
            // ZIP them together before saving the final ZIP to disk.
            downloadControlDocuments(<TDownloadControlDocumentsInput>{
                files: this.selectedFiles,
                orgId: PageParamsStore.state.organization!.Id,
                catId: this.catId,
            }).then((resp : TDownloadControlDocumentsOutput) => {
                saveAs(resp.data, "download.zip")
            }).catch((err : any) => {
                console.log(err)
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        changePage(newPage : number) {
            if (this.pageNumOneIndex == newPage) {
                return
            }
            this.selectedFiles = []
            this.refreshData(newPage - 1, false)
        },

        newRequest(req : DocumentRequest) {
            this.requests.push(req)
            this.showHideRequest = false
        }
    },
    watch : {
        catId() {
            this.refreshData(0, true)
        }
    },
    mounted() {
        this.refreshData(0, true)
    }
})

</script>

<style scoped>

.hidden {
    visibility: hidden;
}

</style>
