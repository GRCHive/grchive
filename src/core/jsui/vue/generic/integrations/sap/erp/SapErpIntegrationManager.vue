<template>
    <div>
        <v-tabs v-if="!loading">
            <v-tab>
                Overview
            </v-tab>
            <v-tab-item>
                <v-row>
                    <v-col cols="6">
                        <v-card>
                            <v-card-title>
                                GRCHive Properties
                            </v-card-title>
                            <v-divider></v-divider>

                            <generic-integration-form
                                class="mx-4 mt-4"
                                v-model="internalIntegration"
                                :valid.sync="genericFormValid"
                                :type="integrationType"
                                :readonly="!genericCanEdit"
                            >
                            </generic-integration-form>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="cancelGenericEdit"
                                    v-if="genericCanEdit"
                                >
                                    Cancel
                                </v-btn>

                                <v-spacer></v-spacer>

                                <v-btn
                                    color="success"
                                    v-if="!genericCanEdit"
                                    @click="genericCanEdit = true"
                                >
                                    Edit
                                </v-btn>

                                <v-btn
                                    color="success"
                                    @click="saveGenericEdit"
                                    :loading="genericIsSaving"
                                    v-else
                                >
                                    Save
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-col>

                    <v-col cols="6">
                        <v-card>
                            <v-card-title>
                                SAP ERP Connection
                            </v-card-title>
                            <v-divider></v-divider>

                            <sap-erp-setup
                                class="mx-4 mt-4"
                                v-model="sapSetup"
                                :valid.sync="sapFormValid"
                                :readonly="!sapCanEdit"
                            >
                            </sap-erp-setup>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="cancelSapEdit"
                                    v-if="sapCanEdit"
                                >
                                    Cancel
                                </v-btn>

                                <v-spacer></v-spacer>

                                <v-btn
                                    color="success"
                                    v-if="!sapCanEdit"
                                    @click="sapCanEdit = true"
                                >
                                    Edit
                                </v-btn>

                                <v-btn
                                    color="success"
                                    @click="saveSapEdit"
                                    :loading="sapIsSaving"
                                    v-else
                                >
                                    Save
                                </v-btn>
                            </v-card-actions>

                        </v-card>
                    </v-col>
                </v-row>
            </v-tab-item>

            <v-tab>
                RFC
            </v-tab>
            <v-tab-item>
                <sap-erp-rfc-manager
                    :integration="integration"
                >
                </sap-erp-rfc-manager>
            </v-tab-item>
        </v-tabs>

        <v-row justify="center" v-else>
            <v-progress-circular size="64" indeterminate></v-progress-circular>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import { GenericIntegration, IntegrationType } from '../../../../../ts/integrations/integration'
import { SapErpIntegrationSetup } from '../../../../../ts/integrations/sap'
import {
    getSapErpIntegration, TGetSapErpIntegrationOutput,
    editSapErpIntegration, TEditSapErpIntegrationOutput
} from '../../../../../ts/api/integrations/apiSapErp'
import {
    editGenericIntegration, TEditGenericIntegrationOutput
} from '../../../../../ts/api/integrations/apiIntegrations'
import { contactUsUrl } from '../../../../../ts/url'
import { PageParamsStore } from '../../../../../ts/pageParams'
import GenericIntegrationForm from '../../GenericIntegrationForm.vue'
import SapErpSetup from './SapErpSetup.vue'
import SapErpRfcManager from './SapErpRfcManager.vue'

@Component({
    components: {
        GenericIntegrationForm,
        SapErpSetup,
        SapErpRfcManager,
    }
})
export default class SapErpIntegrationManager extends Vue {
    @Prop({required: true})
    readonly integration! : GenericIntegration

    internalIntegration : GenericIntegration | null = null
    genericFormValid : boolean = false
    genericCanEdit: boolean = false
    genericIsSaving: boolean = false

    sapSetup : SapErpIntegrationSetup | null = null
    refSap : SapErpIntegrationSetup | null = null
    sapFormValid : boolean = false
    sapCanEdit: boolean = false
    sapIsSaving: boolean = false

    @Watch('integration')
    syncIntegration() {
        this.internalIntegration = JSON.parse(JSON.stringify(this.integration))
    }

    get integrationType() : IntegrationType {
        return IntegrationType.SapErp
    }

    get loading() : boolean {
        return !this.sapSetup
    }

    refreshData() {
        getSapErpIntegration({
            orgId: PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
        }).then((resp : TGetSapErpIntegrationOutput) => {
            this.sapSetup = resp.data
            this.refSap = JSON.parse(JSON.stringify(this.sapSetup!))
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
        this.syncIntegration()
    }

    cancelGenericEdit() {
        this.genericCanEdit = false
        this.syncIntegration()
    }

    saveGenericEdit() {
        this.genericIsSaving = true
        editGenericIntegration({
            orgId: PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
            data: this.internalIntegration!,
        }).then((resp : TEditGenericIntegrationOutput) => {
            this.genericCanEdit = false
            this.$emit('update:integration', resp.data)
            this.internalIntegration = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        }).finally(() => {
            this.genericIsSaving = false
        })
    }

    cancelSapEdit() {
        this.sapCanEdit = false
        this.sapSetup = JSON.parse(JSON.stringify(this.refSap))
    }

    saveSapEdit() {
        this.sapIsSaving = true
        editSapErpIntegration({
            orgId: PageParamsStore.state.organization!.Id,
            integrationId: this.integration.Id,
            setup: this.sapSetup!,
        }).then((resp : TEditSapErpIntegrationOutput) => {
            this.sapCanEdit = false
            this.refSap = resp.data
            this.sapSetup = JSON.parse(JSON.stringify(this.refSap!))
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        }).finally(() => {
            this.sapIsSaving = false
        })
    }
}

</script>
