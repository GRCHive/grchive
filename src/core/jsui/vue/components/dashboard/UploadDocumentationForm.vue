<template>
    <v-card>
        <v-overlay :value="progressOverlay">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>
        <v-card-title>
            Upload Documentation
        </v-card-title>
        <v-divider></v-divider>
        <v-form @submit.prevent v-model="formValid" class="ma-4">
            <document-category-search-form-component
                v-if="catId == -1"
                v-model="chosenCat"
                load-cats
                :rules="[rules.required]"
            >
            </document-category-search-form-component>

            <div v-if="requestId != -1 && requestLinkedToControl && !requestControl && folderId == -1">
                <control-search-form-component
                    v-model="chosenControl"
                    :rules="[rules.required]"
                    :readonly="!!requestControl"
                >
                </control-search-form-component>

                <control-folder-search-form-component
                    :control-id="chosenControl.Id"
                    v-model="chosenFolder"
                    :rules="[rules.required]"
                    v-if="!!chosenControl"
                >
                </control-folder-search-form-component>
            </div>

            <v-menu
                offset-y
                v-model="dateMenu"
                :close-on-content-click="false"
                v-if="!isVersionUpload"
            >
                <template v-slot:activator="{ on }">
                    <v-text-field
                        v-model="dateString"
                        label="Document Date"
                        prepend-icon="mdi-calendar"
                        readonly
                        v-on="on">
                    </v-text-field>
                </template>

                <v-date-picker v-model="dateString" @input="dateMenu = false">
                </v-date-picker>
            </v-menu>

            <div id="drag-drop-file-upload"></div>

            <user-search-form-component
                label="File Owner"
                v-bind:user.sync="uploadUser"
                disabled
            ></user-search-form-component>
        </v-form>

        <v-card-actions>
            <v-btn color="error" @click="onCancel">
                Cancel
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn type="submit" color="primary" :disabled="!canSubmit" @click="submitForm">
                Submit
            </v-btn>
        </v-card-actions>
    </v-card>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import { TUploadControlDocOutput, uploadControlDoc } from '../../../ts/api/apiControlDocumentation'
import { contactUsUrl, uploadControlDocUrl } from '../../../ts/url'
import * as rules from '../../../ts/formRules'
import { PageParamsStore } from '../../../ts/pageParams'
import { createLocalDateFromDateString } from '../../../ts/time'
import { ControlDocumentationCategory } from '../../../ts/controls'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'
import DocumentCategorySearchFormComponent from '../../generic/DocumentCategorySearchFormComponent.vue'
import ControlSearchFormComponent from '../../generic/ControlSearchFormComponent.vue'
import ControlFolderSearchFormComponent from '../../generic/ControlFolderSearchFormComponent.vue'
import { standardFormatDate } from '../../../ts/time'
import { FileFolder } from '../../../ts/folders'
import { getAPIRequestConfig } from '../../../ts/api/apiUtility'

import { Uppy } from '@uppy/core'
import Dashboard from '@uppy/dashboard'
import XHRUpload from '@uppy/xhr-upload'
import '@uppy/core/dist/style.css'
import '@uppy/dashboard/dist/style.css'

@Component({
    components: {
        UserSearchFormComponent,
        DocumentCategorySearchFormComponent,
        ControlSearchFormComponent,
        ControlFolderSearchFormComponent,
    }, 
})
export default class UploadDocumentationForm extends Vue {
    @Prop()
    catId!: number

    @Prop({default: -1})
    requestId!: number

    // If requestId and requestLinkedToControl is set, then that should mean
    // that we are uploading a file to fulfill a request that is linked to
    // a control. If requestControl is set, then that means we know which
    // control to link to and that it shouldn't be changed.
    @Prop({default: false})
    requestLinkedToControl!: boolean

    @Prop({default: null})
    requestControl!: ProcessFlowControl | null

    // If fileId is set, it means we're uploading a new version of the file.
    @Prop({default: -1})
    fileId!: number

    @Prop({default: -1})
    folderId!: number

    dateMenu: boolean = false
    dateString: string = standardFormatDate(new Date())
    rules : any = rules
    formValid: boolean = false
    progressOverlay: boolean = false
    uploadUser: User | null = null
    chosenCat: ControlDocumentationCategory | null = null
    chosenControl: ProcessFlowControl | null = null
    chosenFolder: FileFolder | null = null

    uppy : Uppy | null = null
    uppyFilesDirty: number = 1

    get uppyHasFiles() : boolean {
        return !!this.uppy && (this.uppyFilesDirty > 0) && (this.uppy!.getFiles().length > 0)
    }

    get canSubmit() : boolean {
        return this.formValid && this.finalCatId != -1 && this.uppyHasFiles
    }

    get finalCatId() : number {
        if (this.catId == -1) {
            if (!!this.chosenCat) {
                return this.chosenCat!.Id
            } else {
                return -1
            }
        }
        return this.catId
    }

    get isVersionUpload() : boolean {
        return this.fileId >= 0
    }

    get getRelevantFolderId() : number | null {
        if (this.folderId != -1) {
            return this.folderId
        }

        if (!!this.chosenFolder) {
            return this.chosenFolder.Id
        }

        return null
    }

    clearForm() {
        this.dateString = standardFormatDate(new Date())
        this.chosenCat =  null as ControlDocumentationCategory | null
        this.chosenFolder = null

        this.uploadUser = PageParamsStore.state.user!
        this.chosenControl = this.requestControl

        if (!!this.uppy) { 
            this.uppy.cancelAll()
            this.uppyFilesDirty += 1
        }
    }

    onCancel() {
        this.$emit('do-cancel')
        this.clearForm()
    }

    submitForm() {
        if (!this.canSubmit) {
            return
        }

        this.progressOverlay = true

        // We need to do this conversion since new Date(str) will create the date in UTC.
        // So we need to make sure the date is valid in the current timezone.
        let currentDate : Date = createLocalDateFromDateString(this.dateString)

        this.uppy!.setMeta({
            catId: this.finalCatId,
            orgId: PageParamsStore.state.organization!.Id,
            relevantTime: currentDate.toISOString(),
            uploadUserId: !!this.uploadUser ? this.uploadUser.Id : -1,
        })

        if (this.getRelevantFolderId != null) {
            this.uppy!.setMeta({folderId: this.getRelevantFolderId})
        }

        if (this.requestId != -1) {
            this.uppy!.setMeta({fulfilledRequestId: this.requestId})
        }

        if (this.isVersionUpload) {
            this.uppy!.setMeta({fileId: this.fileId})
        }

        // Need to make sure the "description" metadata is populated with something so the
        // server doesn't bork.
        for (let f of this.uppy!.getFiles()) {
            let m = f.meta
            if (!('description' in m)) {
                this.uppy!.setFileMeta(f.id, {
                    description: ''
                })
            }
        }

        this.uppy!.upload().then(this.finishUpload)
    }

    initUppy() {
        this.uppy = new Uppy({
            debug: false,
            autoProceed: false,
            restrictions: {
                maxNumberOfFiles: this.isVersionUpload ? 1 : null,
                minNumberOfFiles: 1,
            },
            onBeforeFileAdded: (currentFile, files) => {
                const newFile = {
                    ...currentFile,
                    meta: {
                        ...currentFile.meta,
                        altName: currentFile.name,
                        description: '',
                    }
                }
                return newFile
            }
        })
        this.uppy
            .use(Dashboard, {
                target : '#drag-drop-file-upload',
                inline: true,
                waitForThumbnailsBeforeUpload: false,
                showLinkToFileUploadResult: false,
                showProgressDetails: true,
                height: 400,
                hideUploadButton: true,
                hideRetryButton: true,
                hidePauseResumeButton: true,
                hideCancelButton: true,
                proudlyDisplayPoweredByUppy: false,
                metaFields: [
                    { id: 'name', name: 'Name', placeholder:'File Name' },
                    { id: 'altName', name: 'Human-Friendly Name', placeholder:'Alternate Name' },
                    { id: 'description', name: 'Description', placeholder:'Description' },
                ]
            })
            .use(XHRUpload, {
                endpoint: uploadControlDocUrl,
                headers: getAPIRequestConfig().headers,
                limit: 3,
            })
            .on('upload-success', (file, resp) => {
                this.$emit('do-save', resp.body.File, resp.body.Version)
            })
    }

    finishUpload(result : any) {
        for (let f of result.successful) {
            this.uppy!.removeFile(f.id)
        }

        const hasFailed : boolean = result.failed.length > 0
        if (hasFailed) {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);

            this.uppy!.setState({
                currentUploads: {},
                totalProgress: 0,
            })
            this.uppy!.getFiles().forEach(file => {
                this.uppy!.setFileState(file.id, {
                    progress: { uploadComplete: false, uploadStarted: false },
                    error: null,
                })
            })
        }
        this.progressOverlay = false
    }

    mounted() {
        this.clearForm()
        this.initUppy()
    }
}

</script>

<style scoped>

#drag-drop-file-upload {
    margin-bottom: 16px;
}

</style>
