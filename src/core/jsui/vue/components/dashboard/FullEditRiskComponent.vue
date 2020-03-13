<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Risk: {{ fullRiskData.Risk.Name }}
                        <span class="subtitle-1" v-if="fullRiskData.Risk.Name != fullRiskData.Risk.Identifier">
                            ({{ fullRiskData.Risk.Identifier }})
                        </span>
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ fullRiskData.Risk.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-spacer></v-spacer>
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
                            item-name="risk"
                            :items-to-delete="[`${fullRiskData.Risk.Name} (${fullRiskData.Risk.Identifier})`]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="5">
                                <create-new-risk-form 
                                  :node-id="-1"
                                  :edit-mode="true"
                                  :default-name="fullRiskData.Risk.Name"
                                  :default-identifier="fullRiskData.Risk.Identifier"
                                  :default-description="fullRiskData.Risk.Description"
                                  :risk-id="fullRiskData.Risk.Id"
                                  :staged-edits="true"
                                  @do-save="onEditRisk"
                                >
                                </create-new-risk-form>
                            </v-col>

                            <v-col cols="7">
                                <v-card class="mb-4">
                                    <v-card-title>
                                        Related Resources
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <v-tabs>
                                        <v-tab>Process Flows</v-tab>
                                        <v-tab-item>
                                            <process-flow-table
                                                :resources="fullRiskData.Flows"
                                            >
                                            </process-flow-table>
                                        </v-tab-item>

                                        <v-tab>Controls</v-tab>
                                        <v-tab-item>
                                            <control-table
                                                :resources="fullRiskData.Controls"
                                            >
                                            </control-table>
                                        </v-tab-item>

                                        <v-tab>Systems</v-tab>
                                        <v-tab-item>
                                            <systems-table
                                                :resources="relevantSystems"
                                            >
                                            </systems-table>
                                        </v-tab-item>

                                        <v-tab>Accounts</v-tab>
                                        <v-tab-item>
                                            <general-ledger-accounts-table
                                                :resources="relevantAccounts"
                                            >
                                            </general-ledger-accounts-table>
                                        </v-tab-item>

                                    </v-tabs>
                                </v-card>

                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Audit Trail</v-tab>
                    <v-tab-item>
                        <audit-trail-viewer
                            :resource-type="['process_flow_risks']"
                            :resource-id="[`${fullRiskData.Risk.Id}`]"
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
import { FullRiskData } from '../../../ts/risks'
import { getSingleRisk, TSingleRiskInput, TSingleRiskOutput} from '../../../ts/api/apiRisks'
import {
    TAllRiskSystemLinkOutput, allRiskSystemLink
} from '../../../ts/api/apiRiskSystemLinks'
import { createOrgRisksUrl, contactUsUrl } from '../../../ts/url'
import CreateNewRiskForm from './CreateNewRiskForm.vue'
import { System } from '../../../ts/systems'
import SystemsTable from '../../generic/SystemsTable.vue'
import ControlTable from '../../generic/ControlTable.vue'
import ProcessFlowTable from '../../generic/ProcessFlowTable.vue'
import GeneralLedgerAccountsTable from '../../generic/GeneralLedgerAccountsTable.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import { allRiskGLLink, TAllRiskGLLinkOutput } from '../../../ts/api/apiRiskGLLinks'
import { GeneralLedger, GeneralLedgerAccount } from '../../../ts/generalLedger'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'
import { deleteRisk, TDeleteRiskOutput } from '../../../ts/api/apiRisks'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'

export default Vue.extend({
    data: () => ({
        expandDescription: false,
        ready: false,
        showHideDelete: false,
        fullRiskData: Object() as FullRiskData,
        relevantSystems: [] as System[],
        relevantGL: null as GeneralLedger | null
    }),
    computed: {
        relevantAccounts() : GeneralLedgerAccount[] {
            if (!this.relevantGL) {
                return []
            }
            return this.relevantGL.listAccounts
        }
    },
    methods: {
        onEditRisk(risk : ProcessFlowRisk) {
            this.fullRiskData.Risk.Name = risk.Name
            this.fullRiskData.Risk.Description = risk.Description
            this.fullRiskData.Risk.Identifier = risk.Identifier
        },
        refreshSystemLink() {
            allRiskSystemLink({
                riskId: this.fullRiskData.Risk.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllRiskSystemLinkOutput) => {
                this.relevantSystems = resp.data.Systems!
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        refreshGLLink() {
            allRiskGLLink({
                riskId: this.fullRiskData.Risk.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllRiskGLLinkOutput) => {
                this.relevantGL = new GeneralLedger()
                this.relevantGL.rebuildGL(resp.data.Categories!, resp.data.Accounts!)
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        refreshRiskData() {
            let data = window.location.pathname.split('/')
            let riskId = Number(data[data.length - 1])

            getSingleRisk(<TSingleRiskInput>{
                riskId: riskId,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TSingleRiskOutput) => {
                this.fullRiskData = resp.data
                this.ready = true

                this.refreshSystemLink()
                this.refreshGLLink()
            }).catch((err : any) => {
                window.location.replace('/404')
            })
        },
        onDelete() {
            deleteRisk({
                nodeId: -1,
                riskIds: [this.fullRiskData.Risk.Id],
                global: true,
            }).then((resp : TDeleteRiskOutput) => {
                window.location.replace(createOrgRisksUrl(PageParamsStore.state.organization!.OktaGroupName))
            }).catch((err) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }

    },
    components: {
        CreateNewRiskForm,
        SystemsTable,
        ControlTable,
        ProcessFlowTable,
        GeneralLedgerAccountsTable,
        AuditTrailViewer,
        GenericDeleteConfirmationForm,
    },
    mounted() {
        this.refreshRiskData()
    }
})

</script>
