<template>
    <div>
        <v-list-item>
            <v-spacer></v-spacer>
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

export default Vue.extend({
    props : {
        catId: Number,
    },
    data : () => ({
        showHideUpload: false
    }),
    components : {
        UploadDocumentationForm
    },
    methods: {
        cancelUpload() {
            this.showHideUpload = false
        },
        finishUpload(newDoc : any) {
            this.showHideUpload = false
        }
    }
})

</script>
