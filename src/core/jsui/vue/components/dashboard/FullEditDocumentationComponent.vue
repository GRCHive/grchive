<template>
    <div class="max-height">
        <v-row class="max-height">
            <v-col cols="6" v-if="previewReady" class="max-height py-0">
                <pdf-js-viewer :pdf="previewBase64" ref="pdfViewer" :max-viewer-height="viewerMaxHeight">
                </pdf-js-viewer>
            </v-col>

            <v-col 
                cols="6"
                v-else
                align="center"
                justify="center"
            >
                <v-progress-circular
                    indeterminate
                    size="64"
                    v-if="hasPreview"
                ></v-progress-circular>

                <span v-else>No Preview Available</span>
            </v-col>

            <v-col cols="6" v-if="metadataReady" class="py-0">
                <div class="mt-2 mr-3">
                    <v-list-item two-line class="pa-0">
                        <v-list-item-content>
                            <v-list-item-title class="title">
                                File: {{ metadata.AltName }}
                            </v-list-item-title>

                            <v-list-item-subtitle>
                                Parent Category: <a :href="parentCatUrl">{{ parentCat.Name }}</a>
                            </v-list-item-subtitle>

                        </v-list-item-content>

                        <v-list-item-content>
                            <v-select
                                v-model="selectedVersion"
                                :items="versionItems"
                                solo
                                flat
                                hide-details
                                @input="selectVersion"
                            >
                            </v-select>
                        </v-list-item-content>

                        <v-spacer></v-spacer>

                        <v-list-item-action>
                            <v-dialog v-model="showHideDelete"
                                      persistent
                                      max-width="40%"
                            >
                                <template v-slot:activator="{ on }">
                                    <v-btn
                                        color="error"
                                        v-on="on"
                                        :disabled="!versionDataReady"
                                    >
                                        Delete
                                    </v-btn>
                                </template>

                                <generic-delete-confirmation-form
                                    item-name="files"
                                    :items-to-delete="[metadata.AltName]"
                                    :use-global-deletion="false"
                                    @do-cancel="showHideDelete = false"
                                    @do-delete="onDelete">
                                </generic-delete-confirmation-form>
                            </v-dialog>
                        </v-list-item-action>

                        <v-list-item-action>
                            <v-btn
                                color="success"
                                @click="onDownload"
                                :disabled="!versionDataReady"
                                class="ml-4"
                            >
                                Download
                            </v-btn>
                        </v-list-item-action>

                        <v-list-item-action>
                            <v-dialog v-model="showHideUpload"
                                      persistent
                                      max-width="40%"
                            >
                                <template v-slot:activator="{ on }">
                                    <v-btn
                                        color="primary"
                                        v-on="on"
                                    >
                                        Upload
                                    </v-btn>
                                </template>

                                <upload-documentation-form
                                    :cat-id="parentCat.Id"
                                    :file-id="metadata.Id"
                                    @do-cancel="showHideUpload = false"
                                    @do-save="onNewVersion"
                                >
                                </upload-documentation-form>
                            </v-dialog>

                        </v-list-item-action>
                    </v-list-item>

                    <v-divider></v-divider>

                    <v-tabs v-model="currentTab">
                        <v-tab>Overview</v-tab>
                        <v-tab>Comments</v-tab>
                    </v-tabs>

                    <v-tabs-items v-model="currentTab" ref="tabItems">
                        <v-tab-item :style="tabItemStyle">
                            <v-form @submit.prevent v-model="formValid" class="ma-4">

                                <v-text-field v-model="editData.File.AltName"
                                              label="Name"
                                              filled
                                              :rules="[rules.required, rules.createMaxLength(256)]"
                                              :disabled="!canEdit"
                                ></v-text-field>

                                <document-category-search-form-component
                                    v-model="editData.ParentCat"
                                    load-cats
                                    :rules="[rules.required]"
                                    disabled
                                >
                                </document-category-search-form-component>

                                <v-menu
                                    offset-y
                                    v-model="showHideDateMenu"
                                    :close-on-content-click="false"
                                    :disabled="!canEdit"
                                >
                                    <template v-slot:activator="{ on }">
                                        <v-text-field
                                            :value="fileDateString"
                                            label="Document Date"
                                            prepend-icon="mdi-calendar"
                                            readonly
                                            v-on="on">
                                        </v-text-field>
                                    </template>

                                    <v-date-picker :value="fileDateString" @input="changeFileDate">
                                    </v-date-picker>
                                </v-menu>

                                <v-textarea v-model="editData.File.Description"
                                            label="Description"
                                            filled
                                            :disabled="!canEdit"
                                ></v-textarea> 

                            </v-form>

                            <div 
                                class="ml-4 mr-4 mb-4"
                                style="display: flex;"
                            >
                                <v-btn
                                    color="error"
                                    @click="cancelEdit"
                                    v-if="canEdit"
                                >
                                    Cancel
                                </v-btn>
                                <v-spacer></v-spacer>

                                <v-btn
                                    color="primary"
                                    @click="canEdit = true"
                                    v-if="!canEdit"
                                >
                                    Edit
                                </v-btn>

                                <v-btn
                                    color="success"
                                    @click="saveEdit"
                                    v-if="canEdit"
                                >
                                    Save
                                </v-btn>
                            </div>


                            <v-divider class="mb-2"></v-divider>

                            <div
                                class="mx-4"
                                v-if="versionDataReady"
                            >
                                <v-text-field v-model="versionStorageData.StorageName"
                                              label="Filename"
                                              filled
                                              :rules="[rules.required, rules.createMaxLength(256)]"
                                              readonly
                                ></v-text-field>

                                <v-text-field
                                    :value="uploadDateString"
                                    label="Upload Date"
                                    readonly
                                >
                                </v-text-field>

                                <user-search-form-component
                                    label="Upload User"
                                    v-bind:user="uploadUser"
                                    disabled
                                ></user-search-form-component>
                            </div>

                            <v-progress-circular
                                indeterminate
                                size="64"
                                v-else
                            ></v-progress-circular>
                        </v-tab-item>

                        <v-tab-item :style="tabItemStyle">
                            <comment-manager
                                :params="commentParams"
                            ></comment-manager>
                        </v-tab-item>
                    </v-tabs-items>
                </div>
            </v-col>

            <v-col cols="6" v-else align="center" justify="center">
                <v-progress-circular
                    indeterminate
                    size="64"
                ></v-progress-circular>
            </v-col>
        </v-row>
    </div>
</template>

<script lang="ts">

interface EditData {
    ParentCat : ControlDocumentationCategory
    File : ControlDocumentationFile
}

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { contactUsUrl, createSingleDocCatUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { 
    getSingleControlDocument,
    TGetSingleControlDocumentOutput,
    downloadSingleControlDocument,
    TDownloadSingleControlDocumentOutput,
    editControlDoc,
    TEditControlDocOutput,
    deleteControlDocuments,
    getVersionStorageData,
    TGetVersionStorageDataOutput,
} from '../../../ts/api/apiControlDocumentation'
import { 
    ControlDocumentationFile,
    ControlDocumentationCategory,
    cleanJsonControlDocumentationFile,
    FileVersion,
    FileStorageData
} from '../../../ts/controls'
import PdfJsViewer from '../../generic/pdf/PdfJsViewer.vue'
import * as rules from '../../../ts/formRules'
import { createLocalDateFromDateString, standardFormatDate } from '../../../ts/time'
import DocumentCategorySearchFormComponent from '../../generic/DocumentCategorySearchFormComponent.vue'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CommentManager from '../../generic/CommentManager.vue'
import MetadataStore from '../../../ts/metadata'
import UploadDocumentationForm from './UploadDocumentationForm.vue'
import { saveAs } from 'file-saver'

@Component({
    components: {
        PdfJsViewer,
        DocumentCategorySearchFormComponent,
        UserSearchFormComponent,
        GenericDeleteConfirmationForm,
        CommentManager,
        UploadDocumentationForm
    }
})
export default class FullEditDocumentationComponent extends Vue {
    ready: boolean = false

    // Metadata information
    parentCat : ControlDocumentationCategory | null = null
    metadata : ControlDocumentationFile | null = null
    availableVersions : FileVersion[] | null = null
    selectedVersion : FileVersion | null = null

    // Version data
    versionStorageData : FileStorageData | null = null
    uploadUser: User | null = null

    // Preview
    hasPreview: boolean = true
    previewData : Blob | null = null
    previewBase64 : string | null = null


    viewerMaxHeight: number = 100
    metadataMaxHeight:  number = 100

    currentTab : number | null = null

    formValid: boolean = false
    showHideDateMenu: boolean = false
    editData: EditData | null = null
    rules : any = rules
    canEdit: boolean = false

    showHideDelete: boolean = false
    showHideUpload: boolean = false

    $refs!: {
        pdfViewer: PdfJsViewer
        tabItems: any
    }

    get versionItems() : any[] {
        if (!this.availableVersions) {
            return []
        }

        return this.availableVersions.map((ele : FileVersion) => ({
            text: `v${ele.VersionNumber}`,
            value: ele,
        }))
    }

    get commentParams() : Object {
        return {
            fileId: this.metadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }
    }

    get fileDateString() : string {
        if (!this.editData) {
            return ""
        }
        return standardFormatDate(this.editData.File.RelevantTime)
    }

    get uploadDateString() : string {
        if (!this.versionStorageData) {
            return ""
        }
        return standardFormatDate(this.versionStorageData.UploadTime)
    }

    changeFileDate(str : string) {
        this.editData!.File.RelevantTime = createLocalDateFromDateString(str)
        this.showHideDateMenu = false
    }

    get previewReady() : boolean {
        return !!this.previewBase64
    }

    get versionDataReady() : boolean {
        return !!this.versionStorageData
    }

    @Watch('previewData')
    encodePreview() {
        if (!this.previewData) {
            this.previewBase64 = null
            return
        }
        
        let reader = new FileReader()
        reader.onload = (e) => {
            if (!e.target) {
                return
            }

            let blobStr : string = <string>e.target.result
            this.previewBase64 = blobStr.replace(/^data:.*\/.*;base64,/g, '')
        }
        reader.readAsDataURL(this.previewData!)
    }

    get metadataReady() : boolean {
        return !!this.metadata
    }

    get parentCatUrl() : string {
        if (!this.parentCat) {
            return "#"
        }

        return createSingleDocCatUrl(PageParamsStore.state.organization!.OktaGroupName, this.parentCat!.Id)
    }

    reloadPreview() {
        this.previewData = null
        if (!this.hasPreview) {
            return
        }

        downloadSingleControlDocument({
            fileId: this.metadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            version: this.selectedVersion!.VersionNumber,
            preview: true,
        }).then((resp : TDownloadSingleControlDocumentOutput) => {
            this.previewData = resp.data
        }).catch((err : any) => {
            console.log("reload preview: " + err)
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    generateEditData() {
        this.editData = {
            ParentCat: JSON.parse(JSON.stringify(this.parentCat)),
            File: JSON.parse(JSON.stringify(this.metadata)),
        }
        cleanJsonControlDocumentationFile(this.editData.File)
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        let url = new URL(window.location.href)
        let params = url.searchParams
        let versionString : string | null = params.get("version")

        getSingleControlDocument({
            fileId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSingleControlDocumentOutput) => {
            this.parentCat = resp.data.Category
            this.metadata = resp.data.File
            this.availableVersions = resp.data.Versions
            if (this.availableVersions.length > 0) {
                if (!!versionString) {
                    let idx = this.availableVersions.findIndex((ele : FileVersion) => ele.VersionNumber == Number(versionString))
                    if (idx == -1) {
                        this.selectVersion(this.availableVersions[0])
                    } else {
                        this.selectVersion(this.availableVersions[idx])
                    }
                } else {
                    this.selectVersion(this.availableVersions[0])
                }
            }
            this.generateEditData()
        }).catch((err : any) => {
            console.log("refresh data: " + err)
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
        window.addEventListener("resize", this.updateViewerRect)
        window.addEventListener("resize", this.updateMetadataRect)
    }

    @Watch('previewReady')
    updateViewerRect() {
        Vue.nextTick(() => {
            if (!this.previewReady || !this.$refs.pdfViewer) {
                return
            }

            let viewerRect = <DOMRect>this.$refs.pdfViewer.$el.getBoundingClientRect()
            let windowHeight = window.innerHeight
            this.viewerMaxHeight = windowHeight - viewerRect.y
        })
    }

    @Watch('metadataReady')
    updateMetadataRect() {
        Vue.nextTick(() => {
            if (!this.metadataReady) {
                return
            }
            let metadataRect = this.$refs.tabItems.$el.getBoundingClientRect()
            let windowHeight = window.innerHeight
            this.metadataMaxHeight = windowHeight - metadataRect.y
        })
    }

    cancelEdit() {
        this.canEdit = false
        this.generateEditData()
    }

    saveEdit() {
        editControlDoc({
            fileId: this.metadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            relevantTime: this.editData!.File.RelevantTime,
            altName: this.editData!.File.AltName,
            description: this.editData!.File.Description,
        }).then((resp : TEditControlDocOutput) => {
            this.canEdit = false
            this.metadata = resp.data.File
            this.parentCat = resp.data.Category
            this.generateEditData()
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

    get tabItemStyle() : any {
        return {
            "overflow": "auto",
            "height": `${this.metadataMaxHeight}px`,
        }
    }

    onDelete() {
        let fileIds = [this.metadata!.Id]
        deleteControlDocuments({
            fileIds: fileIds,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.showHideDelete = false
            window.location.replace(createSingleDocCatUrl(PageParamsStore.state.organization!.OktaGroupName, this.parentCat!.Id))
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

    onDownload() {
        downloadSingleControlDocument({
            fileId: this.metadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            version: this.selectedVersion!.VersionNumber,
            preview: false,
        }).then((resp : TDownloadSingleControlDocumentOutput) => {
            saveAs(resp.data, this.versionStorageData!.StorageName)
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

    selectVersion(v : FileVersion) {
        this.selectedVersion = v
        this.versionStorageData = null
        getVersionStorageData({
            fileId: this.metadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            version: this.selectedVersion!.VersionNumber,
        }).then((resp : TGetVersionStorageDataOutput) => {
            this.versionStorageData = resp.data.Storage
            this.uploadUser = MetadataStore.getters.getUser(resp.data.Storage.UploadUserId)
            this.hasPreview = resp.data.HasPreview
            this.reloadPreview()
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

    onNewVersion(f : ControlDocumentationFile, v : FileVersion) {
        this.showHideUpload = false
        this.availableVersions!.unshift(v)
        this.selectVersion(v)
    }
}

</script>
