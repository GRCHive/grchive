<template>
    <div>
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
            :resources="value"
        ></doc-file-table>

        <v-divider></v-divider>
        <v-list-item>
            <v-dialog v-model="showHideSelectFiles" persistent max-width="60%">
                <template v-slot:activator="{on}">
                    <v-btn color="warning" v-on="on">
                        Select
                    </v-btn>
                </template>

                <v-card>
                    <v-card-title>
                        Select Files
                    </v-card-title>
                    <v-divider></v-divider>

                    <doc-file-table
                        :resources="value"
                        selectable
                        multi
                        v-model="selectedFiles"
                    ></doc-file-table>

                    <v-card-actions>
                        <v-btn color="error" @click="showHideSelectFiles = false">
                            Cancel
                        </v-btn>

                        <v-spacer></v-spacer>

                        <v-btn color="error" @click="startDeleteFlow" :disabled="!hasSelected || deleteInProgress">
                            <span v-if="!deleteInProgress">Delete</span>
                            <v-progress-circular indeterminate size="16" v-else></v-progress-circular>
                        </v-btn>

                        <slot
                            name="multiActions" 
                            v-bind:hasSelected="hasSelected"
                            v-bind:selectedFiles="selectedFiles"
                        >
                        </slot>

                        <v-btn color="success" @click="downloadSelectedFiles" :disabled="!hasSelected || downloadInProgress">
                            <span v-if="!downloadInProgress">Download</span>
                            <v-progress-circular indeterminate size="16" v-else></v-progress-circular>
                        </v-btn>
                    </v-card-actions>
                </v-card>

            </v-dialog>

            <v-spacer></v-spacer>

            <slot name="singleActions">
            </slot>

            <v-list-item-action v-if="!disableUpload">
                <v-dialog v-model="showHideUpload" persistent max-width="40%">
                    <template v-slot:activator="{on}">
                        <v-btn color="primary" v-on="on">
                            Upload
                        </v-btn>
                    </template>
                    <upload-documentation-form :cat-id="catId"
                        :request-id="requestId"
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
import DocFileTable from './DocFileTable.vue'
import GenericDeleteConfirmationForm from '../components/dashboard/GenericDeleteConfirmationForm.vue'
import UploadDocumentationForm from '../components/dashboard/UploadDocumentationForm.vue'
import { ControlDocumentationFile } from '../../ts/controls'
import { deleteControlDocuments } from '../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import { TDownloadControlDocumentsOutput, downloadControlDocuments } from '../../ts/api/apiControlDocumentation'
import { saveAs } from 'file-saver'

const Props = Vue.extend({
    props: {
        catId: Number,
        requestId: {
            type: Number,
            default: -1
        },
        value: {
            type: Array,
            default: () => [] as ControlDocumentationFile[]
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
    },
    components: {
        DocFileTable,
        GenericDeleteConfirmationForm,
        UploadDocumentationForm
    }
})

// A wrapper around DocFileTable that provides commonly
// needed controls related to managing files (mainly
// upload/delete/sampling).
@Component
export default class DocFileManager extends Props {
    selectedFiles: ControlDocumentationFile[] = []
    showHideSelectFiles: boolean = false

    showHideDeleteFiles: boolean = false
    showHideUpload: boolean = false

    deleteInProgress : boolean = false
    downloadInProgress : boolean = false

    get canSelect() : boolean {
        return this.forceEnableSelect || !this.disableDownload || !this.disableDelete
    }

    get hasSelected() : boolean {
        return this.selectedFiles.length > 0
    }

    get selectedFileNames() : string[] {
        return this.selectedFiles.map((ele) => ele.StorageName)
    }

    deleteSelectedFiles() {
        this.showHideSelectFiles = true
        this.showHideDeleteFiles = false
        this.deleteInProgress = true
        deleteControlDocuments({
            orgId: PageParamsStore.state.organization!.Id,
            fileIds: this.selectedFiles.map((ele) => ele.Id),
            catId: this.catId,
        }).then(() => {
            let selectedFileSet = new Set(this.selectedFiles)
            for (let i = this.value.length - 1; i >= 0; --i) {
                if (selectedFileSet.has(this.value[i] as ControlDocumentationFile)) {
                    this.value.splice(i, 1)
                    this.$emit('input', this.value)
                }
            }
            this.selectedFiles = []
            this.showHideDeleteFiles = false
            this.deleteInProgress = false
            this.showHideSelectFiles = false
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

    downloadSelectedFiles() {
        this.downloadInProgress = true
        // Download each file individually from the webserver and then
        // ZIP them together before saving the final ZIP to disk.
        downloadControlDocuments({
            files: this.selectedFiles,
            orgId: PageParamsStore.state.organization!.Id,
            catId: this.catId,
        }).then((resp : TDownloadControlDocumentsOutput) => {
            this.downloadInProgress = false
            this.showHideSelectFiles = false
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

        newDoc.RelevantTime = new Date(newDoc.RelevantTime)
        newDoc.UploadTime = new Date(newDoc.RelevantTime)
        this.value.unshift(newDoc)
        this.$emit('new-doc', newDoc)
        this.$emit('input', this.value)
    }

    startDeleteFlow() {
        this.showHideSelectFiles = false
        this.showHideDeleteFiles = true
    }
}

</script>
