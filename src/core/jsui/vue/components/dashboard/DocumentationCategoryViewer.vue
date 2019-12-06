<template>
    <div>
        <v-list-item>
            <v-list-item-action class="ma-1 hidden">
                <v-checkbox>
                </v-checkbox>
            </v-list-item-action>

            <v-list-item-content>
                <v-list-item-title class="font-weight-bold">
                    Filename
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-content>
                <v-list-item-title class="font-weight-bold">
                    Relevant Date
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-content>
                <v-list-item-title class="font-weight-bold">
                    Upload Date
                </v-list-item-title>
            </v-list-item-content>
        </v-list-item>
        <v-divider></v-divider>

        <v-list>
            <v-list-item-group multiple v-model="selectedFiles">
                <v-list-item v-for="(item, index) in files"
                             :key="index"
                             :value="item"
                >
                    <template v-slot:default="{active, toggle}">
                        <v-list-item-action class="ma-1">
                            <v-checkbox :input-value="active"
                                        @true-value="item"
                                        @click="toggle">
                            </v-checkbox>
                        </v-list-item-action>

                        <v-list-item-content>
                            <v-list-item-title>
                                {{index+1}}.&nbsp;{{ item.StorageName }}
                            </v-list-item-title>
                        </v-list-item-content>

                        <v-list-item-content>
                            <v-list-item-title>
                                {{ item.RelevantTime.toDateString() }}
                            </v-list-item-title>
                        </v-list-item-content>

                        <v-list-item-content>
                            <v-list-item-title>
                                {{ item.UploadTime.toDateString() }}
                            </v-list-item-title>
                        </v-list-item-content>
                    </template>
                </v-list-item>
            </v-list-item-group>
        </v-list>

        <div class="text-center">
            <v-pagination
                :value="pageNumOneIndex"
                :length="totalPages"
                @input="changePage"
            ></v-pagination>
        </div>

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
            <v-btn color="success" @click="downloadSelectedFiles" :disabled="!hasSelected">
                Download
            </v-btn>

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
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { saveAs } from 'file-saver'
import { PageParamsStore } from '../../../ts/pageParams'

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
        showHideDeleteFiles: false,
    }),
    components : {
        UploadDocumentationForm,
        GenericDeleteConfirmationForm
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
                fileIds: this.selectedFiles.map((ele) => ele.Id)
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
                orgId: PageParamsStore.state.organization!.Id
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
        },
        changePage(newPage : number) {
            if (this.pageNumOneIndex == newPage) {
                return
            }
            this.selectedFiles = []
            this.refreshData(newPage - 1, false)
        },
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
