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

            <v-col cols="6" v-if="metadataReady">
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

                        <v-spacer></v-spacer>

                        <v-list-item-action>
                            <v-btn
                                color="error"
                            >
                                Delete
                            </v-btn>
                        </v-list-item-action>

                        <v-list-item-action>
                            <v-btn
                                color="success"
                            >
                                Download
                            </v-btn>
                        </v-list-item-action>
                    </v-list-item>

                    <v-divider></v-divider>

                    <v-tabs v-model="currentTab">
                        <v-tab>Overview</v-tab>
                        <v-tab>Comments</v-tab>
                    </v-tabs>

                    <v-tabs-items v-model="currentTab">
                        <v-tab-item style="overflow: auto;">
                            <v-form @submit.prevent v-model="formValid" class="ma-4">

                                <v-text-field v-model="editData.File.AltName"
                                              label="Name"
                                              filled
                                              :rules="[rules.required, rules.createMaxLength(256)]"
                                              :disabled="!canEdit"
                                ></v-text-field>

                                <v-text-field v-model="editData.File.StorageName"
                                              label="Filename"
                                              filled
                                              :rules="[rules.required, rules.createMaxLength(256)]"
                                              readonly
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

                                <v-text-field
                                    :value="uploadDateString"
                                    label="Upload Date"
                                    readonly
                                >
                                </v-text-field>

                                <v-textarea v-model="editData.File.Description"
                                            label="Description"
                                            filled
                                            :disabled="!canEdit"
                                ></v-textarea> 

                                <user-search-form-component
                                    label="File Owner"
                                    v-bind:user.sync="editData.UploadUser"
                                    :disabled="!canEdit"
                                ></user-search-form-component>
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
                        </v-tab-item>

                        <v-tab-item>
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
    UploadUser : User
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
    TEditControlDocOutput
} from '../../../ts/api/apiControlDocumentation'
import { ControlDocumentationFile, ControlDocumentationCategory, cleanJsonControlDocumentationFile } from '../../../ts/controls'
import PdfJsViewer from '../../generic/pdf/PdfJsViewer.vue'
import * as rules from '../../../ts/formRules'
import { createLocalDateFromDateString, standardFormatDate } from '../../../ts/time'
import DocumentCategorySearchFormComponent from '../../generic/DocumentCategorySearchFormComponent.vue'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'

@Component({
    components: {
        PdfJsViewer,
        DocumentCategorySearchFormComponent,
        UserSearchFormComponent
    }
})
export default class FullEditDocumentationComponent extends Vue {
    ready: boolean = false
    parentCat : ControlDocumentationCategory | null = null
    metadata : ControlDocumentationFile | null = null
    uploadUser: User | null = null
    viewerMaxHeight: number = 100

    previewMetadata : ControlDocumentationFile | null = null
    previewData : Blob | null = null
    previewBase64 : string | null = null

    currentTab : number | null = null

    formValid: boolean = false
    showHideDateMenu: boolean = false
    editData: EditData | null = null
    rules : any = rules
    canEdit: boolean = false

    $refs!: {
        pdfViewer: PdfJsViewer
    }

    get fileDateString() : string {
        if (!this.editData) {
            return ""
        }
        return standardFormatDate(this.editData.File.RelevantTime)
    }

    get uploadDateString() : string {
        if (!this.editData) {
            return ""
        }
        return standardFormatDate(this.editData.File.UploadTime)
    }

    changeFileDate(str : string) {
        this.editData!.File.RelevantTime = createLocalDateFromDateString(str)
        this.showHideDateMenu = false
    }

    get hasPreview() : boolean {
        return !!this.previewMetadata
    }

    get previewReady() : boolean {
        return !!this.previewBase64
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
        if (!this.previewMetadata) {
            return
        }

        downloadSingleControlDocument({
            fileId: this.previewMetadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            catId: this.previewMetadata!.CategoryId,
        }).then((resp : TDownloadSingleControlDocumentOutput) => {
            this.previewData = resp.data
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

    generateEditData() {
        this.editData = {
            ParentCat: JSON.parse(JSON.stringify(this.parentCat)),
            File: JSON.parse(JSON.stringify(this.metadata)),
            UploadUser: JSON.parse(JSON.stringify(this.uploadUser))
        }
        cleanJsonControlDocumentationFile(this.editData.File)
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getSingleControlDocument({
            fileId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSingleControlDocumentOutput) => {
            this.parentCat = resp.data.Category
            this.previewMetadata = resp.data.PreviewFile
            this.uploadUser = resp.data.UploadUser
            this.metadata = resp.data.File
            this.generateEditData()
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

    mounted() {
        this.refreshData()
        window.addEventListener("resize", this.updateViewerRect)
    }

    @Watch('previewReady')
    updateViewerRect() {
        Vue.nextTick(() => {
            if (!this.previewReady || !this.$refs.pdfViewer) {
                return
            }

            let viewerRect = this.$refs.pdfViewer.$el.getBoundingClientRect()
            let windowHeight = window.innerHeight
            this.viewerMaxHeight = windowHeight - viewerRect.y
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
            uploadUserId: this.editData!.UploadUser.Id
        }).then((resp : TEditControlDocCatOutput) => {
            this.canEdit = false
            this.metadata = resp.data.File
            this.parentCat = resp.data.Category
            this.uploadUser = resp.data.UploadUser
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
}

</script>
