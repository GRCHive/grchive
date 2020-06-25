<template>
    <div>
        <v-list-item>
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    {{ customHeader != "" ? customHeader : "Files" }}
                </v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <v-dialog v-model="showHideDeleteFiles" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="documents"
                :items-to-delete="selectedFileNames"
                v-on:do-cancel="showHideDeleteFiles = false"
                v-on:do-delete="deleteSelectedFiles"
                :use-global-deletion="false"
                :force-global-deletion="true">
            </generic-delete-confirmation-form>
        </v-dialog>

        <doc-file-table
            v-model="selectedFiles"
            :resources="allFiles"
            :search="filterText"
            :selectable="selectMode"
            :multi="selectMode"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            :value="value"
            @delete="deleteSingleFile"
            @input="modifySelected"
        ></doc-file-table>

        <v-divider></v-divider>
        <v-list-item>
            <v-list-item-action class="mr-2">
                <v-btn
                    color="warning"
                    @click="selectMode = true"
                    v-if="!selectMode"
                >
                    Select
                </v-btn>

                <v-btn
                    color="error"
                    @click="selectMode = false"
                    v-if="selectMode"
                >
                    Cancel
                </v-btn>
            </v-list-item-action>

            <v-list-item-action class="mr-2" v-if="!disableDelete">
                <v-btn color="error" @click="startDeleteFlow" :disabled="!hasSelected" :loading="deleteInProgress">
                    <span v-if="!deleteInProgress">Delete</span>
                    <v-progress-circular indeterminate size="16" v-else></v-progress-circular>
                </v-btn>
            </v-list-item-action>

            <slot
                name="multiActions" 
                v-bind:hasSelected="hasSelected"
                v-bind:selectedFiles="selectedFiles"
                v-bind:allFiles="allFiles"
            >
            </slot>

            <v-list-item-action class="mr-2" v-if="!disableDownload">
                <v-btn color="success" @click="downloadSelectedFiles" :disabled="!hasSelected" :loading="downloadInProgress">
                    <span v-if="!downloadInProgress">Download</span>
                    <v-progress-circular indeterminate size="16" v-else></v-progress-circular>
                </v-btn>
            </v-list-item-action>

            <v-spacer></v-spacer>

            <slot name="singleActions">
            </slot>

            <v-list-item-action v-if="canLinkFiles">
                <v-dialog v-model="showHideAddExisting" persistent max-width="40%">
                    <template v-slot:activator="{on}">
                        <v-btn color="secondary" v-on="on">
                            Add Existing
                        </v-btn>
                    </template>

                    <doc-searcher-form
                        :exclude-files="allFiles"
                        :force-cat-id="catId"
                        @do-cancel="showHideAddExisting = false"
                        @do-select="doLinkFiles"
                    >
                    </doc-searcher-form>
                </v-dialog>

            </v-list-item-action>

            <v-list-item-action v-if="!disableUpload">
                <v-dialog v-model="showHideUpload" persistent max-width="40%">
                    <template v-slot:activator="{on}">
                        <v-btn color="primary" v-on="on">
                            Upload
                        </v-btn>
                    </template>
                    <upload-documentation-form
                        :cat-id="catId"
                        :folder-id="folderId"
                        :request-id="requestId"
                        :request-linked-to-control="requestLinkedToControl"
                        :request-control="requestControl"
                        @do-cancel="showHideUpload = false"
                        @do-save="finishUpload">
                    </upload-documentation-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import DocFileTable from './DocFileTable.vue'
import DocSearcherForm from './DocSearcherForm.vue'
import GenericDeleteConfirmationForm from '../components/dashboard/GenericDeleteConfirmationForm.vue'
import UploadDocumentationForm from '../components/dashboard/UploadDocumentationForm.vue'
import { ControlDocumentationFile, VersionedMetadata } from '../../ts/controls'
import { deleteControlDocuments } from '../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import { 
    TDownloadControlDocumentsOutput, downloadControlDocuments,
    allControlDocuments, TAllControlDocumentsOutput,
} from '../../ts/api/apiControlDocumentation'
import {
    cleanJsonControlDocumentationFile
} from '../../ts/controls'
import {
    allFolderFileLink, TAllFolderFileLinkOutput ,
    newFolderFileLink,
    deleteFolderFileLink
} from '../../ts/api/apiFolderFileLinks'
import { 
    getVendorProduct, TGetVendorProductOutput,
    linkVendorProductSocFiles, unlinkVendorProductSocFiles,
} from '../../ts/api/apiVendorProduct'
import {
    TGetSingleDocumentRequestOutput, getSingleDocRequest,
    TLinkFilesToDocumentRequestInput, linkFilesToDocRequest,
} from '../../ts/api/apiDocRequests'
import { extractControlDocumentationFileHandle } from '../../ts/controls'
import { FileFolder } from '../../ts/folders'

import { saveAs } from 'file-saver'

const Props = Vue.extend({
    props: {
        catId: {
            type: Number,
            default: -1,
        },
        folder: {
            type: Object,
            default: () => null as FileFolder | null
        },
        customHeader: {
            type: String,
            default: "",
        },
        vendorProductId: {
            type: Number,
            default: -1,
        },
        vendorId: {
            type: Number,
            default: -1,
        },
        requestId: {
            type: Number,
            default: -1
        },
        requestLinkedToControl: {
            type: Boolean,
            default: false,
        },
        requestControl: {
            type: Object,
            default: () => null as ProcessFlowControl | null
        },
        disableUpload: {
            type: Boolean,
            default: false
        },
        disableDelete: {
            type: Boolean,
            default: false
        },
        disableSample: {
            type: Boolean,
            default: false
        },
        disableDownload: {
            type: Boolean,
            default: false
        },
        forceEnableSelect: {
            type: Boolean,
            default: true,
        },
        value: {
            type: Array,
            default: () => [],
        },

    },
    components: {
        DocFileTable,
        GenericDeleteConfirmationForm,
        UploadDocumentationForm,
        DocSearcherForm
    }
})

// A wrapper around DocFileTable that provides commonly
// needed controls related to managing files (mainly
// upload/delete/sampling).
@Component
export default class DocFileManager extends Props {
    allFiles : ControlDocumentationFile[] = []
    filterText: string = ""

    selectedFiles: VersionedMetadata[] = []
    selectMode: boolean = false

    showHideDeleteFiles: boolean = false
    showHideUpload: boolean = false
    showHideAddExisting : boolean = false

    deleteInProgress : boolean = false
    downloadInProgress : boolean = false

    get folderId() : number {
        if (!!this.folder) {
            return this.folder.Id
        }
        return -1
    }

    get canSelect() : boolean {
        return this.forceEnableSelect || !this.disableDownload || !this.disableDelete
    }

    get isFolderMode() : boolean {
        return !!this.folder
    }

    get isVendorProductSocMode() : boolean {
        return this.vendorId != -1 && this.vendorProductId != -1
    }

    get isDocRequestMode() : boolean {
        return this.requestId != -1
    }

    get canLinkFiles() : boolean {
        return !this.disableUpload && (this.isFolderMode || this.isVendorProductSocMode)
    }

    get hasSelected() : boolean {
        return this.selectedFiles.length > 0
    }

    get selectedFileNames() : string[] {
        return this.selectedFiles.map((ele) => ele.File.AltName)
    }

    downloadSelectedFiles() {
        this.downloadInProgress = true
        // Download each file individually from the webserver and then
        // ZIP them together before saving the final ZIP to disk.
        downloadControlDocuments({
            files: this.selectedFiles,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TDownloadControlDocumentsOutput) => {
            this.downloadInProgress = false
            saveAs(resp.data, "download.zip")
        }).catch((err : any) => {
            this.downloadInProgress = false
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    finishUpload(newDoc : ControlDocumentationFile) {
        this.showHideUpload = false
        cleanJsonControlDocumentationFile(newDoc)
        this.allFiles.unshift(newDoc)
        this.$emit('new-doc', newDoc)
    }

    startDeleteFlow() {
        this.showHideDeleteFiles = true
    }

    @Watch('catId')
    @Watch('folder')
    @Watch('vendorId')
    @Watch('vendorProductId')
    @Watch('requestId')
    refreshData() {
         if (this.isDocRequestMode) {
            getSingleDocRequest({
                requestId: this.requestId,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TGetSingleDocumentRequestOutput) => {
                this.allFiles = resp.data.Files
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        } else if (this.isVendorProductSocMode) {
            getVendorProduct({
                productId: this.vendorProductId,
                vendorId: this.vendorId,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TGetVendorProductOutput) => {
                this.allFiles = resp.data.SocFiles
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        } else if (this.isFolderMode) {
            allFolderFileLink({
                folderId: this.folder!.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllFolderFileLinkOutput) => {
                this.allFiles = resp.data.Files!
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        } else if (this.catId != -1) {
            allControlDocuments({
                catId: this.catId,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllControlDocumentsOutput) => {
                this.allFiles = resp.data.Files
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

    mounted() {
        this.refreshData()
    }

    genericDeleteFiles(ids : number[], clearSelect:  boolean) {
        let onSuccess = (sids : number[]) => {
            let selectedFileSet = new Set<number>(sids)
            for (let i = this.allFiles.length - 1; i >= 0; --i) {
                if (selectedFileSet.has((this.allFiles[i] as ControlDocumentationFile).Id)) {
                    this.allFiles.splice(i, 1)
                }
            }
            if (clearSelect) {
                this.selectedFiles = []
            }
            this.showHideDeleteFiles = false
            this.deleteInProgress = false
        }

        if (this.isVendorProductSocMode) {
            unlinkVendorProductSocFiles({
                productId: this.vendorProductId,
                vendorId: this.vendorId,
                orgId: PageParamsStore.state.organization!.Id,
                socFiles: ids.map((ele : number) => ({
                    Id: ele,
                    CategoryId: this.catId,
                }))
            }).then(() => {
                onSuccess(ids)
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        } else if (this.isFolderMode) {
            for (let id of ids) {
                deleteFolderFileLink({
                    folderId: this.folder!.Id,
                    fileId: id,
                    orgId: PageParamsStore.state.organization!.Id,
                }).then(() => {
                    onSuccess([id])
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
        } else {
            deleteControlDocuments({
                orgId: PageParamsStore.state.organization!.Id,
                fileIds: ids,
            }).then(() => {
                onSuccess(ids)
            }).catch((err : any) => {
                this.deleteInProgress = false
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

    deleteSelectedFiles() {
        this.showHideDeleteFiles = false
        this.deleteInProgress = true

        let selectedFileIds = this.selectedFiles.map((ele :VersionedMetadata) => ele.File.Id)
        this.genericDeleteFiles(selectedFileIds, true)
    }

    deleteSingleFile(file : ControlDocumentationFile) {
        this.genericDeleteFiles([file.Id], false)
    }

    modifySelected(vals : VersionedMetadata[]) {
        this.$emit('input', vals)
    }

    doLinkFiles(files : ControlDocumentationFile[]) {
        if (this.isDocRequestMode) {
            let params : TLinkFilesToDocumentRequestInput = {
                requestId: this.requestId,
                orgId: PageParamsStore.state.organization!.Id,
                files: files.map((ele : ControlDocumentationFile) => ele.Id),
            }

            linkFilesToDocRequest(params).then(() => {
                this.showHideAddExisting = false
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

        if (this.isVendorProductSocMode) {
            linkVendorProductSocFiles({
                productId: this.vendorProductId,
                vendorId: this.vendorId,
                orgId: PageParamsStore.state.organization!.Id,
                socFiles: files.map(extractControlDocumentationFileHandle)
            }).then(() => {
                this.showHideAddExisting = false
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

        if (this.isFolderMode) {
            newFolderFileLink({
                folderId: this.folder!.Id,
                fileIds: files.map((ele : ControlDocumentationFile) => ele.Id),
                orgId: PageParamsStore.state.organization!.Id,
            }).then(() => {
                this.showHideAddExisting = false
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

        this.allFiles.unshift(...files)
    }
}

</script>
