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
import { getCurrentCSRF } from '../../../ts/csrf'

export default Vue.extend({
    props : {
        catId: Number
    },
    data : () => ({
        dateMenu: false,
        dateString: new Date().toISOString().substr(0, 10),
        file: null as File | null,
        rules,
        formValid: false,
        progressOverlay: false
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

            let data = new FormData()
            data.set('csrf', getCurrentCSRF())
            data.set('catId', this.catId.toString())
            data.set('file', this.file!)
            data.set('relevantTime', currentDate.toISOString())

            uploadControlDoc(data).then((resp : TUploadControlDocOutput) => {
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
})

</script>

