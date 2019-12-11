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
            <v-menu offset-y v-model="dateMenu" :close-on-content-click="false">
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
            ></v-text-field>

            <v-textarea v-model="description"
                        label="Description"
                        filled
            ></v-textarea> 

            <user-search-form-component
                label="File Owner"
                v-bind:user.sync="uploadUser"
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
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'

export default Vue.extend({
    props : {
        catId: Number,
        requestId: {
            type: Number,
            default: -1
        }
    },
    components: {
        UserSearchFormComponent
    },
    data : () => ({
        dateMenu: false,
        dateString: new Date().toISOString().substr(0, 10),
        file: null as File | null,
        rules,
        formValid: false,
        progressOverlay: false,
        altName: "",
        description: "",
        uploadUser: {} as User
    }),
    computed : {
        canSubmit() : boolean {
            return !!this.file && this.formValid
        },
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
            // We need to do this conversion so that when we go into UTC, we're still at the right
            // day for the local time zone.
            let currentDate : Date = new Date(this.dateString)
            currentDate = new Date(currentDate.getTime() + currentDate.getTimezoneOffset() * 60000)

            uploadControlDoc({
                catId: this.catId,
                orgId: PageParamsStore.state.organization!.Id,
                file: this.file!,
                relevantTime: currentDate,
                altName: this.altName,
                description: this.description,
                uploadUserId: this.uploadUser.Id,
                fulfilledRequestId: (this.requestId != -1) ? this.requestId : null
            }).then((resp : TUploadControlDocOutput) => {
                this.progressOverlay = false
                this.$emit('do-save', resp.data)
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
        console.log(this.uploadUser)
    }
})

</script>

