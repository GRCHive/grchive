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
import { standardFormatDate } from '../../../ts/time'

export default Vue.extend({
    props : {
        catId: Number,
        requestId: {
            type: Number,
            default: -1
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
        DocumentCategorySearchFormComponent
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
        chosenCat: {} as ControlDocumentationCategory,
    }),
    computed : {
        canSubmit() : boolean {
            return !!this.file && this.formValid
        },

        finalCatId() : number {
            if (this.catId == -1) {
                return this.chosenCat.Id
            }
            return this.catId
        },

        isVersionUpload() : boolean {
            return this.fileId >= 0
        }
    },
    methods: {
        onCancel() {
            this.$emit('do-cancel')
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
                folderId: (this.folderId != -1) ? this.folderId : null,
            }).then((resp : TUploadControlDocOutput) => {
                this.progressOverlay = false
                this.$emit('do-save', resp.data.File, resp.data.Version)
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
        this.uploadUser = PageParamsStore.state.user!
    }
})

</script>
