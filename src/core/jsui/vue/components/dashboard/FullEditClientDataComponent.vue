<template>
    <div>
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line>
                <v-list-item-content id="reqTitle">
                    <v-list-item-title class="title">
                        Data Object: {{ data.Data.Name }}
                    </v-list-item-title>

                    <v-list-item-subtitle v-if="!!source">
                        Source:
                        <resource-handle-renderer :handle="source"></resource-handle-renderer>
                    </v-list-item-subtitle>
                </v-list-item-content>
                <v-spacer></v-spacer>

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
                        item-name="data objects"
                        :items-to-delete="[data.Data.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid class="pa-0">
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="12">
                                <create-new-client-data-form
                                    edit-mode
                                    :reference-data="data"
                                    @do-save="onEditData"
                                    class="ma-4"
                                >
                                </create-new-client-data-form>
                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Code</v-tab>
                    <v-tab-item>
                        <v-divider></v-divider>
                        <managed-code-ide
                            :data-id="data.Data.Id"
                            lang="text/x-kotlin"
                            full-height
                        >
                        </managed-code-ide>

                    </v-tab-item>

                    <v-tab>Audit Trail</v-tab>
                    <v-tab-item>
                        <audit-trail-viewer
                            :resource-type="['client_data']"
                            :resource-id="[`${data.Data.Id}`]"
                            no-header
                        >
                        </audit-trail-viewer>

                    </v-tab-item>
                </v-tabs>
            </v-container>
        </div>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { FullClientDataWithLink } from '../../../ts/clientData'
import { ResourceHandle } from '../../../ts/resourceUtils'
import { getDataSource, TGetDataSourceOutput } from '../../../ts/api/apiDataSource'
import { getClientData, TGetClientDataOutput, deleteClientData } from '../../../ts/api/apiClientData'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgClientDataUrl } from '../../../ts/url'
import { 
    allCode, TAllCodeOutput,
} from '../../../ts/api/apiCode'
import { ManagedCode } from '../../../ts/code'

import ResourceHandleRenderer from '../../generic/ResourceHandleRenderer.vue'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'
import CreateNewClientDataForm from './CreateNewClientDataForm.vue'
import ManagedCodeIde from '../../generic/code/ManagedCodeIDE.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'

@Component({
    components: {
        ResourceHandleRenderer,
        AuditTrailViewer,
        CreateNewClientDataForm,
        ManagedCodeIde,
        GenericDeleteConfirmationForm,
    }
})
export default class FullEditClientDataComponent extends Vue {
    data : FullClientDataWithLink | null = null
    source : ResourceHandle | null = null
    showHideDelete : boolean = false

    get ready() : boolean {
        return !!this.data
    }

    onEditData(data : FullClientDataWithLink) {
        this.data = data
        this.refreshSource()
    }

    refreshSource() {
        this.source = null
        getDataSource({
            source: this.data!.Link
        }).then((resp : TGetDataSourceOutput) => {
            this.source = resp.data
        }).catch((err: any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    refreshData() {
        this.data = null
        this.source = null

        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getClientData({
            orgId: PageParamsStore.state.organization!.Id,
            dataId: resourceId,
        }).then((resp : TGetClientDataOutput) => {
            this.data = resp.data
            this.refreshSource()
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

    onDelete() {
        deleteClientData({
            orgId: PageParamsStore.state.organization!.Id,
            dataId: this.data!.Data.Id,
        }).then(() => {
            window.location.replace(createOrgClientDataUrl(PageParamsStore.state.organization!.OktaGroupName))
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
