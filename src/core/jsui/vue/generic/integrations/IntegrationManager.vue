<template>
    <div>
        <div v-if="!loading">
            <v-list-item>
                <v-list-item-action>
                    <v-dialog persistent max-width="40%" v-model="showHideDelete">
                        <template v-slot:activator="{on}">
                            <v-btn
                                color="error"
                                :disabled="!hasCurrentIntegration"
                                v-on="on"
                            >
                                Delete
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            v-if="!!currentIntegration"
                            item-name="integrations"
                            :items-to-delete="[currentIntegration.Name]"
                            @do-cancel="showHideDelete = false"
                            @do-delete="deleteCurrentIntegration"
                            :use-global-deletion="false"
                            :delete-in-progress="deleteInProgress"
                        >
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

                <v-spacer></v-spacer>

                <v-list-item-action>
                    <v-dialog
                        v-model="showHideChooser"
                        persistent
                        max-width="60%"
                    >
                        <template v-slot:activator="{on}">
                            <v-btn
                                color="primary"
                                v-on="on"
                            >
                                New
                            </v-btn>
                        </template>

                        <v-card>
                            <integration-new-setup
                                ref="setup"
                                @cancel="showHideChooser = false"
                                @save="onNewIntegration"
                            ></integration-new-setup>
                        </v-card>
                    </v-dialog>
                </v-list-item-action>
            </v-list-item>

            <v-tabs vertical v-model="currentTab">
                <template v-for="(data, index) in storedIntegrations">
                    <v-tab>
                        {{ data.Name }}
                    </v-tab>

                    <v-tab-item>
                        <sap-erp-integration-manager
                            v-if="data.Type == IntegrationType.SapErp"
                            :integration.sync="storedIntegrations[index]"
                        >
                        </sap-erp-integration-manager>
                    </v-tab-item>
                </template>
            </v-tabs>
        </div>

        <v-row justify="center" v-else>
            <v-progress-circular size="64" indeterminate></v-progress-circular>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import IntegrationNewSetup from './IntegrationNewSetup.vue'
import { contactUsUrl } from '../../../ts/url'
import { IntegrationType, GenericIntegration, UnionIntegrationSetup } from '../../../ts/integrations/integration'
import { SapErpIntegrationSetup } from '../../../ts/integrations/sap'
import {
    newSapErpIntegration, TNewSapErpIntegrationOutput
} from '../../../ts/api/integrations/apiSapErp'
import {
    allGenericIntegrations, TAllGenericIntegrationsOutput,
    deleteGenericIntegration,
} from '../../../ts/api/integrations/apiIntegrations'
import { TIntegrationLink } from '../../../ts/api/integrations/apiIntegrations'
import { PageParamsStore } from '../../../ts/pageParams'
import GenericDeleteConfirmationForm from '../../components/dashboard/GenericDeleteConfirmationForm.vue'

const SapErpIntegrationManager = () => import (/* webpackChunkName: "SapErpIntegrationManager" */ './sap/erp/SapErpIntegrationManager.vue')

@Component({
    components: {
        SapErpIntegrationManager,
        IntegrationNewSetup,
        GenericDeleteConfirmationForm,
    }
})
export default class IntegrationManager extends Vue {
    @Prop()
    readonly systemId : number | undefined

    showHideChooser: boolean = false
    showHideDelete: boolean = false
    storedIntegrations : GenericIntegration[] | null = null
    IntegrationType : any = IntegrationType
    deleteInProgress:  boolean = false

    currentTab : number = 0

    $refs! : {
        setup: IntegrationNewSetup
    }

    get loading() : boolean {
        return !this.storedIntegrations
    }

    get integrationLink() : TIntegrationLink {
        return {
            systemId: this.systemId
        }
    }

    get hasCurrentIntegration():  boolean {
        return (this.currentTab >= 0 && this.currentTab < this.storedIntegrations!.length)
    }

    get currentIntegration() : GenericIntegration | null {
        if (!this.hasCurrentIntegration) {
            return null
        }

        return this.storedIntegrations![this.currentTab]
    }

    deleteCurrentIntegration() {
        if (!this.currentIntegration) {
            return
        }

        let id = this.currentIntegration.Id
        this.deleteInProgress = true

        deleteGenericIntegration({
            orgId: PageParamsStore.state.organization!.Id,
            integrationId: id,
        }).then(() => {
            let idx = this.storedIntegrations!.findIndex((ele : GenericIntegration) => ele.Id == id)
            if (idx == -1) {
                return
            }
            this.storedIntegrations!.splice(idx, 1)
            this.showHideDelete = false
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

    onNewIntegration(integration : GenericIntegration, setup : UnionIntegrationSetup) {
        if (!setup) {
            this.onError()
            return
        }

        switch (integration.Type) {
        case IntegrationType.SapErp:
            this.onNewSapErpIntegration(integration, <SapErpIntegrationSetup>setup)
            break
        default:
            break
        }
    }

    onNewSapErpIntegration(integration : GenericIntegration, sapSetup : SapErpIntegrationSetup) {
        newSapErpIntegration({
            orgId: PageParamsStore.state.organization!.Id,
            integration: integration,
            setup: sapSetup,
            link: this.integrationLink,
        }).then((resp : TNewSapErpIntegrationOutput) => {
            this.storedIntegrations!.push(resp.data)
            this.onSuccess()
        }).catch(this.onError)
    }

    onError() {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops! Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true)

        if (!!this.$refs.setup) {
            this.$refs.setup.clearSaving()
        }
    }

    onSuccess() {
        this.showHideChooser = false
        this.$refs.setup.cleanup()
    }

    refreshData() {
        allGenericIntegrations({
            orgId: PageParamsStore.state.organization!.Id,
            link: this.integrationLink,
        }).then((resp : TAllGenericIntegrationsOutput) => {
            this.storedIntegrations = resp.data
        }).catch(this.onError)
    }

    mounted() {
        this.refreshData()
    }
}

</script>
