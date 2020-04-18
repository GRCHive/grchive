<template>
    <div>
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line>
                <v-list-item-content id="reqTitle">
                    <v-list-item-title class="title">
                        Script: {{ data.Name }}
                    </v-list-item-title>
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
                        item-name="script"
                        :items-to-delete="[data.Name]"
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
                                <create-new-script-form
                                    class="ma-4"
                                    edit-mode
                                    :reference-script="data"
                                    @do-save="onEditScript"
                                >
                                </create-new-script-form>
                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Code</v-tab>
                    <v-tab-item>
                        <v-divider></v-divider>
                        <managed-code-ide
                            :script-id="data.Id"
                            lang="text/x-kotlin"
                            full-height
                        >
                        </managed-code-ide>
                    </v-tab-item>

                    <v-tab>Run Logs</v-tab>
                    <v-tab-item>
                        <run-log-list
                            :script-id="data.Id"
                        >
                        </run-log-list>
                    </v-tab-item>

                    <v-tab>Audit Trail</v-tab>
                    <v-tab-item>
                    </v-tab-item>
                </v-tabs>
            </v-container>
        </div>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgScriptUrl } from '../../../ts/url'
import { ClientScript } from '../../../ts/clientScripts'
import {
    getClientScript, TGetClientScriptOutput,
    deleteClientScript
} from '../../../ts/api/apiScripts'
import CreateNewScriptForm from './CreateNewScriptForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import ManagedCodeIde from '../../generic/code/ManagedCodeIDE.vue'
import RunLogList from '../../generic/logs/RunLogList.vue'

@Component({
    components: {
        CreateNewScriptForm,
        GenericDeleteConfirmationForm,
        ManagedCodeIde,
        RunLogList
    }
})
export default class FullEditClientScriptComponent extends Vue {
    data : ClientScript | null = null
    showHideDelete : boolean = false

    get ready() : boolean {
        return !!this.data
    }

    onEditScript(data : ClientScript) {
        this.data = data
    }

    refreshData() {
        this.data = null

        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getClientScript({
            orgId: PageParamsStore.state.organization!.Id,
            scriptId: resourceId,
        }).then((resp : TGetClientScriptOutput) => {
            this.data = resp.data
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
        deleteClientScript({
            orgId: PageParamsStore.state.organization!.Id,
            scriptId: this.data!.Id,
        }).then(() => {
            window.location.replace(createOrgScriptUrl(PageParamsStore.state.organization!.OktaGroupName))
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
