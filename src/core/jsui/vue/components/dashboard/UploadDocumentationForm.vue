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

            <div v-if="requestId != -1 && requestLinkedToControl">
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
            <v-file-input label="Documentation" v-model="file" :rules="[rules.required]"></v-file-input>

            <v-text-field v-model="altName"
                          label="Name"
                          filled
                          :rules="[rules.required, rules.createMaxLength(256)]"
                          v-if="!isVersionUpload"
            ></v-text-field>

            <v-textarea v-model="description"
                        label="Description"
                        filled
                        v-if="!isVersionUpload"
            ></v-textarea> 

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
import { TUploadControlDocOutput, uploadControlDoc } from '../../../ts/api/apiControlDocumentation'
import { contactUsUrl } from '../../../ts/url'
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

export default Vue.extend({
    props : {
        catId: Number,
        requestId: {
            type: Number,
            default: -1
        },
        // If requestId and requestLinkedToControl is set, then that should mean
        // that we are uploading a file to fulfill a request that is linked to
        // a control. If requestControl is set, then that means we know which
        // control to link to and that it shouldn't be changed.
        requestLinkedToControl: {
            type: Boolean,
            default: false,
        },
        requestControl: {
            type: Object,
            default: () => null as ProcessFlowControl | null
        },
        // If fileId is set, it means we're uploading a new version of the file.
        fileId:{
            type: Number,
            default: -1
        },
        folderId: {
            type: Number,
            default: -1
        }
    },
    components: {
        UserSearchFormComponent,
        DocumentCategorySearchFormComponent,
        ControlSearchFormComponent,
        ControlFolderSearchFormComponent,
    },
    data : () => ({
        dateMenu: false,
        dateString: standardFormatDate(new Date()),
        file: null as File | null,
        rules,
        formValid: false,
        progressOverlay: false,
        altName: "",
        description: "",
        uploadUser: {} as User,
        chosenCat: null as ControlDocumentationCategory | null,
        chosenControl: null as ProcessFlowControl | null,
        chosenFolder: null as FileFolder | null,
    }),
    computed : {
        canSubmit() : boolean {
            return !!this.file && this.formValid && this.finalCatId != -1
        },

        finalCatId() : number {
            if (this.catId == -1) {
                if (!!this.chosenCat) {
                    return this.chosenCat!.Id
                } else {
                    return -1
                }
            }
            return this.catId
        },

        isVersionUpload() : boolean {
            return this.fileId >= 0
        },

        getRelevantFolderId() : number | null{
            if (this.folderId != -1) {
                return this.folderId
            }

            if (!!this.chosenFolder) {
                return this.chosenFolder.Id
            }

            return null
        }
    },
    methods: {
        clearForm() {
            this.dateString = standardFormatDate(new Date())
            this.file = null
            this.altName = ""
            this.description =  ""
            this.chosenCat =  null as ControlDocumentationCategory | null
            this.chosenFolder = null

            this.uploadUser = PageParamsStore.state.user!
            this.chosenControl = this.requestControl
        },
        onCancel() {
            this.$emit('do-cancel')
            this.clearForm()
        },
        submitForm() {
            if (!this.canSubmit) {
                return
            }
            this.progressOverlay = true
            // We need to do this conversion since new Date(str) will create the date in UTC.
            // So we need to make sure the date is valid in the current timezone.
            let currentDate : Date = createLocalDateFromDateString(this.dateString)

            uploadControlDoc({
                catId: this.finalCatId,
                orgId: PageParamsStore.state.organization!.Id,
                file: this.file!,
                relevantTime: currentDate,
                altName: this.altName,
                description: this.description,
                uploadUserId: this.uploadUser.Id,
                fulfilledRequestId: (this.requestId != -1) ? this.requestId : null,
                fileId: this.isVersionUpload ? this.fileId : null,
                folderId: this.getRelevantFolderId,
            }).then((resp : TUploadControlDocOutput) => {
                this.progressOverlay = false
                this.$emit('do-save', resp.data.File, resp.data.Version)
                this.clearForm()
            }).catch((err : any) => {
                this.progressOverlay = false
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    },
    watch: {
        file() {
            if (!!this.file) {
                this.altName = this.file!.name
            }
        }
    },
    mounted() {
        this.clearForm()
    }
})

</script>
