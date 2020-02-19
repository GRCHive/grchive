<template>
    <div class="ma-4">
        <v-dialog v-model="showHideDeleteRisk" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="risks"
                :items-to-delete="currentRisksToDelete"
                v-on:do-cancel="showHideDeleteRisk = false"
                v-on:do-delete="deleteSelectedRisks"
                :use-global-deletion="true"
                :force-global-deletion="true">
            </generic-delete-confirmation-form>
        </v-dialog>

        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Risks
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
                <v-dialog v-model="showHideCreateNewRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            New
                        </v-btn>
                    </template>
                    <create-new-risk-form
                        :node-id="-1"
                        @do-save="saveNewRisk"
                        @do-cancel="cancelNewRisk">
                    </create-new-risk-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <risk-table
            :resources="allRisks"
            :search="filterText"
            use-crud-delete
            confirm-delete
            @delete="deleteSelectedRisk"
        >
        </risk-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { getAllRisks, TAllRiskInput, TAllRiskOutput } from '../../../ts/api/apiRisks'
import { deleteRisk, TDeleteRiskInput, TDeleteRiskOutput } from '../../../ts/api/apiRisks'
import { contactUsUrl } from '../../../ts/url'
import CreateNewRiskForm from './CreateNewRiskForm.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import RiskTable from '../../generic/RiskTable.vue'

export default Vue.extend({
    data : () => ({
        allRisks: [] as ProcessFlowRisk[],
        filterText : "",
        showHideCreateNewRisk: false,
    }),
    components: {
        CreateNewRiskForm,
        RiskTable
    },
    methods: {
        refreshRisks() {
            getAllRisks(<TAllRiskInput>{
                orgName: PageParamsStore.state.organization!.OktaGroupName
            }).then((resp : TAllRiskOutput) => {
                this.allRisks = resp.data
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        saveNewRisk(risk : ProcessFlowRisk) {
            this.allRisks.unshift(risk)
            this.showHideCreateNewRisk = false
        },
        cancelNewRisk() {
            this.showHideCreateNewRisk = false
        },
        deleteSelectedRisk(risk : ProcessFlowRisk, global : boolean) {
            deleteRisk(<TDeleteRiskInput>{
                nodeId: -1,
                riskIds: [risk.Id],
                global: global,
            }).then((resp : TDeleteRiskOutput) => {
                this.allRisks.splice(
                    this.allRisks.findIndex((ele : ProcessFlowRisk) =>
                        ele.Id == risk.Id),
                    1)
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
    mounted() {
        this.refreshRisks()
    }
})
</script>
