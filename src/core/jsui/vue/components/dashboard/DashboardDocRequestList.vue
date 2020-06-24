<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    PBC Requests
                </v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-dialog v-model="showHideNew"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <create-new-request-form
                        :available-cats="availableCats"
                        @do-cancel="showHideNew = false"
                        @do-save="newRequest">
                    </create-new-request-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <doc-request-table
            :resources="requests"
            :search="filterText"
        ></doc-request-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DocRequestTable from '../../generic/DocRequestTable.vue'
import { DocumentRequest } from '../../../ts/docRequests'
import { ControlDocumentationCategory } from '../../../ts/controls'
import { TGetAllDocumentRequestOutput, getAllDocRequests } from '../../../ts/api/apiDocRequests'
import { TGetAllDocumentationCategoriesOutput, getAllDocumentationCategories } from '../../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import CreateNewRequestForm from './CreateNewRequestForm.vue'

@Component({
    components: {
        DocRequestTable,
        CreateNewRequestForm
    }
})
export default class DashboardDocRequestList extends Vue {
    showHideNew: boolean = false
    filterText: string = ""
    requests : DocumentRequest[] = []
    availableCats: ControlDocumentationCategory[] = []

    mounted() {
        getAllDocRequests({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetAllDocumentRequestOutput) => {
            this.requests = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })

        getAllDocumentationCategories({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetAllDocumentationCategoriesOutput) => {
            this.availableCats = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    newRequest(req : DocumentRequest) {
        this.showHideNew = false
        this.requests.unshift(req)
    }
}

</script>
