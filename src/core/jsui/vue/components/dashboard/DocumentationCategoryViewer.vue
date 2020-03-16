<template>
    <div>
        <doc-file-manager
            class="mb-4"
            :cat-id="catId"
        >
        </doc-file-manager>
        <v-card>
            <v-card-title>
                <span class="mr-2">
                    Requests
                </span>
                <v-spacer></v-spacer>

                <v-dialog v-model="showHideRequest" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="warning" icon v-on="on">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <create-new-request-form
                        :cat-id="catId"
                        @do-cancel="showHideRequest = false"
                        @do-save="newRequest">
                    </create-new-request-form>
                </v-dialog>
            </v-card-title>
            <v-divider></v-divider>

            <doc-request-table :resources="requests">
            </doc-request-table>
        </v-card>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { ControlDocumentationFile } from '../../../ts/controls'
import { allControlDocuments, TAllControlDocumentsOutput } from '../../../ts/api/apiControlDocumentation'
import { getAllDocRequests, TGetAllDocumentRequestOutput } from '../../../ts/api/apiDocRequests'
import { PageParamsStore } from '../../../ts/pageParams'
import DocFileManager from '../../generic/DocFileManager.vue'
import CreateNewRequestForm from './CreateNewRequestForm.vue'
import { DocumentRequest } from '../../../ts/docRequests'
import DocRequestTable from '../../generic/DocRequestTable.vue'
import { contactUsUrl } from '../../../ts/url'

const Props = Vue.extend({
    props : {
        catId: Number,
    }
})


@Component({
    components: {
        DocFileManager,
        CreateNewRequestForm,
        DocRequestTable
    }
})
export default class DocumentationCategoryViewer extends Props {
    requests : DocumentRequest[] = []
    showHideRequest : boolean = false
    
    refreshData() {
        getAllDocRequests({
            orgId: PageParamsStore.state.organization!.Id,
            catId: this.catId,
        }).then((resp : TGetAllDocumentRequestOutput) => {
            this.requests = resp.data
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

    newRequest(req : DocumentRequest) {
        this.requests.push(req)
        this.showHideRequest = false
    }

    mounted() {
        this.refreshData()
    }
}

</script>
