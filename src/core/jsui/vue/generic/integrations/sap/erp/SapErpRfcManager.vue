<template>
    <div>
        <v-row v-if="!loading">
            <v-col cols="12">
                <v-select
                    v-model="currentRfc"
                    label="RFC"
                    hide-details
                    :items="rfcItems"
                    filled
                    dense
                >
                    <template v-slot:append-outer>
                        <v-dialog persistent max-width="40%" v-model="showHideNew">
                            <template v-slot:activator="{on}">
                                <v-btn
                                    icon
                                    color="primary"
                                    v-on="on"
                                >
                                    <v-icon>mdi-plus</v-icon>
                                </v-btn>
                            </template>

                            <v-card>
                                <v-card-title>New RFC</v-card-title>
                                <v-divider></v-divider>

                                <v-form v-model="newRfcValid" class="ma-4">
                                    <v-text-field
                                        label="Function"
                                        v-model="newRfcFunction"
                                        filled
                                        :rules="[rules.required, rules.createMaxLength(512)]"
                                    >
                                    </v-text-field>
                                </v-form>

                                <v-card-actions>
                                    <v-btn
                                        color="error"
                                        @click="showHideNew = false"
                                    >
                                        Cancel
                                    </v-btn>

                                    <v-spacer></v-spacer>

                                    <v-btn
                                        color="success"
                                        :disabled="!newRfcValid"
                                        @click="saveNewRfc"
                                        :loading="newRfcSaving"
                                    >
                                        Save
                                    </v-btn>
                                </v-card-actions>
                            </v-card>
                        </v-dialog>

                        <v-dialog persistent max-width="40%" v-model="showHideDelete">
                            <template v-slot:activator="{on}">
                                <v-btn
                                    icon
                                    color="error"
                                    :disabled="!currentRfc"
                                    v-on="on"
                                >
                                    <v-icon>mdi-delete</v-icon>
                                </v-btn>
                            </template>

                            <generic-delete-confirmation-form
                                v-if="!!currentRfc"
                                item-name="SAP ERP RFCs"
                                :items-to-delete="[currentRfc.Function]"
                                @do-cancel="showHideDelete = false"
                                @do-delete="deleteCurrentRfc"
                                :use-global-deletion="false"
                                :delete-in-progress="deleteInProgress"
                            >
                            </generic-delete-confirmation-form>

                        </v-dialog>
                    </template>
                </v-select>

                <v-select
                    v-model="currentVersion"
                    label="Version"
                    hide-details
                    :items="rfcVersionItems"
                    filled
                    v-if="!!availableVersions"
                    :loading="loadingVersions"
                    dense
                    class="mt-4"
                >
                    <template v-slot:append-outer>
                        <v-btn
                            icon
                            color="success"
                            @click="createRfcVersion"
                            :loading="newVersionInProgress"
                        >
                            <v-icon>mdi-refresh</v-icon>
                        </v-btn>
                    </template>
                </v-select>
                <v-divider class="my-4"></v-divider>

                <sap-erp-rfc-version-data-viewer
                    :integration-id="integration.Id"
                    :version="currentVersion"
                    v-if="!!currentVersion"
                >
                </sap-erp-rfc-version-data-viewer>
            </v-col>
        </v-row>

        <v-row justify="center" v-else>
            <v-progress-circular size="64" indeterminate></v-progress-circular>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import { GenericIntegration } from '../../../../../ts/integrations/integration'
import { contactUsUrl } from '../../../../../ts/url'
import { PageParamsStore } from '../../../../../ts/pageParams'
import {
    SapErpRfcMetadata,
    SapErpRfcVersion
} from '../../../../../ts/integrations/sap'
import {
    allSapErpRfc, TAllSapErpRfcOutput,
    newSapErpRfc, TNewSapErpRfcOutput,
    deleteSapErpRfc,
    allSapErpRfcVersions, TAllSapErpRfcVersionsOutput,
    newSapErpRfcVersion, TNewSapErpRfcVersionOutput,
} from '../../../../../ts/api/integrations/apiSapErp'
import GenericDeleteConfirmationForm from '../../../../components/dashboard/GenericDeleteConfirmationForm.vue'
import SapErpRfcVersionDataViewer from './SapErpRfcVersionDataViewer.vue'
import * as rules from '../../../../../ts/formRules'
import {standardFormatTime} from '../../../../../ts/time'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        SapErpRfcVersionDataViewer,
    }
})
export default class SapErpRfcManager extends Vue {
    @Prop({required: true})
    readonly integration! : GenericIntegration

    rules: any = rules

    showHideNew : boolean = false
    newRfcValid : boolean = false
    newRfcFunction: string = ""
    newRfcSaving : boolean = false

    showHideDelete : boolean = false
    deleteInProgress : boolean = false

    newVersionInProgress : boolean = false

    rfcs : SapErpRfcMetadata[] | null = null
    availableVersions: SapErpRfcVersion[] | null = null
    loadingVersions : boolean = false

    currentRfc : SapErpRfcMetadata | null = null
    currentVersion : SapErpRfcVersion | null = null

    get rfcItems() : any[] {
        return this.rfcs!.map((ele : SapErpRfcMetadata) => ({
            text: ele.Function,
            value: ele,
        }))
    }

    get rfcVersionItems() : any[] {
        if (!this.availableVersions) {
            return []
        }
        return this.availableVersions.map((ele : SapErpRfcVersion, idx : number) => ({
            text: `${this.availableVersions!.length - idx}. ${standardFormatTime(ele.CreatedTime)}`,
            value: ele,
        }))
    }

    get loading() : boolean {
        return !this.rfcs
    }

    @Watch('currentRfc')
    refreshVersions() {
        if (!this.currentRfc) {
            this.availableVersions = null
            this.currentVersion = null
            return
        }

        this.loadingVersions = true
        allSapErpRfcVersions({
            orgId : PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
            rfcId: this.currentRfc!.Id,
        }).then((resp : TAllSapErpRfcVersionsOutput) => {
            this.availableVersions = resp.data
            if (this.availableVersions!.length > 0) {
                this.currentVersion = this.availableVersions![0]
            }
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        }).finally(() => {
            this.loadingVersions = false
        })
    }

    createRfcVersion() {
        this.newVersionInProgress = true
        newSapErpRfcVersion({
            orgId : PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
            rfcId: this.currentRfc!.Id,
        }).then((resp : TNewSapErpRfcVersionOutput) => {
            this.availableVersions!.unshift(resp.data)
            this.currentVersion = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        }).finally(() => {
            this.newVersionInProgress = false
        })
    }

    deleteCurrentRfc() {
        this.deleteInProgress = true
        deleteSapErpRfc({
            orgId : PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
            rfcId: this.currentRfc!.Id,
        }).then(() => {
            this.showHideDelete = false
            let idx : number = this.rfcs!.findIndex((ele : SapErpRfcMetadata) => ele.Id == this.currentRfc!.Id)
            if (idx == -1) {
                return
            }
            this.rfcs!.splice(idx, 1)
            if (this.rfcs!.length > 0) {
                this.currentRfc = this.rfcs![0]
            } else {
                this.currentRfc = null
            }
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        }).finally(() => {
            this.deleteInProgress = false
        })
    }

    saveNewRfc() {
        this.newRfcSaving = true
        newSapErpRfc({
            orgId : PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
            function: this.newRfcFunction,
        }).then((resp : TNewSapErpRfcOutput) => {
            this.newRfcFunction = ""
            this.showHideNew = false
            this.rfcs!.unshift(resp.data)
            this.currentRfc = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        }).finally(() => {
            this.newRfcSaving = false
        })
    }

    refreshData() {
        allSapErpRfc({
            orgId : PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
        }).then((resp : TAllSapErpRfcOutput) => {
            this.rfcs = resp.data
            if (this.rfcs!.length > 0) {
                this.currentRfc = this.rfcs![0]
            }
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>
