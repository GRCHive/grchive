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
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ fullRiskData.Risk.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="5">
                        <create-new-risk-form ref="editRisk"
                                              :node-id="-1"
                                              :edit-mode="true"
                                              :default-name="fullRiskData.Risk.Name"
                                              :default-description="fullRiskData.Risk.Description"
                                              :risk-id="fullRiskData.Risk.Id"
                                              :staged-edits="true"
                                              @do-save="onEditRisk">
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
                            </v-tabs>
                        </v-card>

                    </v-col>
                </v-row>
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
import { createControlUrl, contactUsUrl } from '../../../ts/url'
import CreateNewRiskForm from './CreateNewRiskForm.vue'
import { System } from '../../../ts/systems'
import SystemsTable from '../../generic/SystemsTable.vue'
import ControlTable from '../../generic/ControlTable.vue'
import ProcessFlowTable from '../../generic/ProcessFlowTable.vue'
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    data: () => ({
        expandDescription: false,
        ready: false,
        fullRiskData: Object() as FullRiskData,
        relevantSystems: [] as System[]
    }),
    methods: {
        onEditRisk(risk : ProcessFlowRisk) {
            this.fullRiskData.Risk.Name = risk.Name
            this.fullRiskData.Risk.Description = risk.Description

            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editRisk.clearForm()
            })
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
        refreshRiskData() {
            let data = window.location.pathname.split('/')
            let riskId = Number(data[data.length - 1])

            getSingleRisk(<TSingleRiskInput>{
                riskId: riskId,
            }).then((resp : TSingleRiskOutput) => {
                this.fullRiskData = resp.data
                this.ready = true

                this.refreshSystemLink()

                Vue.nextTick(() => {
                    //@ts-ignore
                    this.$refs.editRisk.clearForm()
                })
            }).catch((err : any) => {
                window.location.replace('/404')
            })
        },
        generateControlUrl(controlId : number) : string {
            return createControlUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                controlId)
        }
    },
    components: {
        CreateNewRiskForm,
        SystemsTable,
        ControlTable,
        ProcessFlowTable
    },
    mounted() {
        this.refreshRiskData()
    }
})

</script>
