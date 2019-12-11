<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Document Request: {{ currentRequest.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentRequest.Description }}
                    </v-list-item-subtitle>

                </v-list-item-content>

                <v-list-item-action>
                    <v-dialog v-model="showHideDelete"
                              persistent
                              max-width="40%"
                    >
                        <template v-slot:activator="{ on }">
                            <v-btn color="error" v-on="on">
                                Delete
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            item-name="document requests"
                            :items-to-delete="[currentRequest.Name]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>

                </v-list-item-action>

                <v-list-item-action>
                    <v-btn color="success" v-if="!currentRequest.CompletionTime">
                        Complete  
                    </v-btn>

                    <v-btn color="error" v-else>
                        Reopen
                    </v-btn>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="8">
                    </v-col>

                    <v-col cols="4">
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { DocumentRequest } from '../../../ts/docRequests'
import { TGetSingleDocumentRequestOutput, getSingleDocRequest } from '../../../ts/api/apiDocRequests'
import { deleteSingleDocRequest } from '../../../ts/api/apiDocRequests'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgDocRequestsUrl } from '../../../ts/url'
import { ControlDocumentationCategory, ControlDocumentationFile } from '../../../ts/controls'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'

@Component({
    components: {
        GenericDeleteConfirmationForm,
    }
})
export default class FullEditDatabaseComponent extends Vue {
    currentRequest : DocumentRequest | null = null
    relevantFiles: ControlDocumentationFile[] = []
    parentCategory : ControlDocumentationCategory | null = null
    expandDescription: boolean = false
    showHideDelete: boolean = false

    get ready() : boolean {
        return !!this.currentRequest && !!this.parentCategory
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getSingleDocRequest({
            requestId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSingleDocumentRequestOutput) => {
            this.currentRequest = resp.data.Request
            this.relevantFiles = resp.data.Files
            this.parentCategory = resp.data.Category
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

    onDelete() {
        deleteSingleDocRequest({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            window.location.replace(createOrgDocRequestsUrl(PageParamsStore.state.organization!.OktaGroupName))
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
    }
}

</script>
