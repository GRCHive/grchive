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
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="8">
                                <create-new-client-data-form
                                    edit-mode
                                    :reference-data="data"
                                    @do-save="onEditData"
                                >
                                </create-new-client-data-form>
                            </v-col>

                            <v-col cols="4">
                                <v-card class="mb-4">
                                    <v-card-title>
                                        Linked Scripts
                                    </v-card-title>
                                    <v-divider></v-divider>
                                </v-card>
                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Code</v-tab>
                    <v-tab-item>
                        <generic-code-editor
                            lang="text/x-kotlin"
                            full-height
                        >
                        </generic-code-editor>
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
import { getClientData, TGetClientDataOutput} from '../../../ts/api/apiClientData'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'

import ResourceHandleRenderer from '../../generic/ResourceHandleRenderer.vue'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'
import CreateNewClientDataForm from './CreateNewClientDataForm.vue'
import GenericCodeEditor from '../../generic/code/GenericCodeEditor.vue'

@Component({
    components: {
        ResourceHandleRenderer,
        AuditTrailViewer,
        CreateNewClientDataForm,
        GenericCodeEditor,
    }
})
export default class FullEditClientDataComponent extends Vue {
    data : FullClientDataWithLink | null = null
    source : ResourceHandle | null = null

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
}

</script>
