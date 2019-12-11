<template>
    <div>
        <doc-file-table
            :resources="value"
            :selectable="!disableDelete"
            :multi="!disableDelete"
            v-model="selectedFiles"
        ></doc-file-table>

        <v-divider></v-divider>
        <v-list-item>
            <span v-if="!disableDelete">
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
            </span>

            <v-spacer></v-spacer>

            <v-list-item-action v-if="!disableSample">
                <v-btn color="info">
                    Sample
                </v-btn>
            </v-list-item-action>

            <v-list-item-action class="ml-4" v-if="!disableDownload">
                <v-btn color="success" @click="downloadSelectedFiles" :disabled="!hasSelected">
                    Download
                </v-btn>
            </v-list-item-action>

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
    showHideDeleteFiles: boolean = false
    showHideUpload: boolean = false

    get hasSelected() : boolean {
        return this.selectedFiles.length > 0
    }

    get selectedFileNames() : string[] {
        return this.selectedFiles.map((ele) => ele.StorageName)
    }

    toggleSelection() {
        if (this.hasSelected) {
            this.selectedFiles = [] as ControlDocumentationFile[]
        } else {
            this.selectedFiles = this.value.slice() as ControlDocumentationFile[]
        }
    }

    deleteSelectedFiles() {
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

    downloadSelectedFiles() {
        // Download each file individually from the webserver and then
        // ZIP them together before saving the final ZIP to disk.
        downloadControlDocuments({
            files: this.selectedFiles,
            orgId: PageParamsStore.state.organization!.Id,
            catId: this.catId,
        }).then((resp : TDownloadControlDocumentsOutput) => {
            saveAs(resp.data, "download.zip")
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

    finishUpload(newDoc : ControlDocumentationFile) {
        this.showHideUpload = false

        newDoc.RelevantTime = new Date(newDoc.RelevantTime)
        newDoc.UploadTime = new Date(newDoc.RelevantTime)
        this.value.unshift(newDoc)
        this.$emit('input', this.value)
    }
}

</script>
