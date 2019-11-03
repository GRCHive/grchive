<template>
    <section>
        <v-dialog v-model="showHideEditRisk" persistent max-width="40%">
            <create-new-risk-form
                ref="editRisk"
                :edit-mode="true"
                :node-id="currentNode.Id"
                :risk-id="editRiskData.Id"
                :default-name="editRiskData.Name"
                :default-description="editRiskData.Description"
                @do-save="finishEditRisk"
                @do-cancel="cancelEditRisk">
            </create-new-risk-form>
        </v-dialog>

        <v-list-item class="pa-1">
            <v-list-item-action class="ma-1">
                <v-btn icon @click="toggleSelection">
                    <v-icon v-if="!hasSelected">mdi-checkbox-blank-outline</v-icon>
                    <v-icon v-else>mdi-minus-box-outline</v-icon>
                </v-btn>
            </v-list-item-action>
            <div class="flex-grow-1"></div>
            <v-list-item-action class="ma-1">
                <v-dialog v-model="showHideDeleteRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="error" :disabled="!hasSelected" v-on="on">
                            Delete
                            <v-icon small>mdi-delete</v-icon>
                        </v-btn>
                    </template>
                    
                    <generic-delete-confirmation-form
                        item-name="risks"
                        :items-to-delete="risksForDeletionConfirmation"
                        v-on:do-cancel="showHideDeleteRisk = false"
                        v-on:do-delete="deleteSelectedRisks"
                        :use-global-deletion="true">
                    </generic-delete-confirmation-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

        <v-list two-line>
            <v-list-item-group multiple v-model="selectedRisks">
                <section v-for="(item, index) in risksForNode"
                         :key="index">
                    <v-list-item :key="index" class="pa-1" :value="item">
                        <template v-slot:default="{active, toggle}">
                            <v-list-item-action class="ma-1">
                                <v-checkbox :input-value="active"
                                            @true-value="item"
                                            @click="toggle">
                                </v-checkbox>
                            </v-list-item-action>

                            <v-list-item-content>
                                <v-list-item-title>
                                    {{ item.Name }}
                                </v-list-item-title>

                                <v-list-item-subtitle>
                                    {{ item.Description }}
                                </v-list-item-subtitle>
                            </v-list-item-content>

                            <v-list-item-action>
                                <v-btn icon @click.stop="editRisk(item)" @mousedown.stop>
                                    <v-icon>mdi-pencil</v-icon>
                                </v-btn>
                            </v-list-item-action>

                            <v-list-item-action>
                                <v-btn icon @click.stop @mousedown.stop :href="generateRiskUrl(item)" target="_blank">
                                    <v-icon>mdi-open-in-new</v-icon>
                                </v-btn>
                            </v-list-item-action>

                        </template>
                    </v-list-item>
                    <v-divider></v-divider>
                </section>
            </v-list-item-group>
        </v-list>

        <v-list-item class="pa-1">
            <v-list-item-action class="ma-1">
                <v-dialog v-model="showHideAddExistingRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            Add Existing
                        </v-btn>
                    </template>
                    <generic-add-existing-item-form
                        item-name="Risks"
                        :selectable-items="unselectedRisksForNode"
                        @do-select="addExistingRisk"
                        @do-cancel="cancelAddRisk">
                    </generic-add-existing-item-form>
                </v-dialog>
            </v-list-item-action>
            <v-list-item-action class="ma-1">
                <v-dialog v-model="showHideCreateNewRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="success" v-on="on">
                            Create New
                        </v-btn>
                    </template>
                    <create-new-risk-form
                        :node-id="currentNode.Id"
                        @do-save="saveNewRisk"
                        @do-cancel="cancelNewRisk">
                    </create-new-risk-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import CreateNewRiskForm from './CreateNewRiskForm.vue'
import DeleteRiskForm from './DeleteRiskForm.vue'
import GenericAddExistingItemForm from './GenericAddExistingItemForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { deleteRisk, addExistingRisk } from '../../../ts/api/apiRisks'
import { contactUsUrl, createRiskUrl } from '../../../ts/url'
import { getCurrentCSRF } from '../../../ts/csrf'

export default Vue.extend({
    data : () => ({
        showHideDeleteRisk : false,
        showHideCreateNewRisk : false,
        showHideAddExistingRisk : false,
        showHideEditRisk : false,
        selectedRisks : [] as ProcessFlowRisk[],
        editRiskData: <ProcessFlowRisk>{
            Id: -1,
            Name: "",
            Description: ""
        }
    }),
    components : {
        CreateNewRiskForm,
        GenericDeleteConfirmationForm,
        GenericAddExistingItemForm
    },
    computed : {
        hasSelected() : boolean {
            return this.selectedRisks.length > 0
        },
        currentNode() : ProcessFlowNode {
            return VueSetup.store.getters.currentNodeInfo
        },
        risksForNode() : ProcessFlowRisk[] {
            return VueSetup.store.getters.risksForNode(this.currentNode.Id)
        },
        unselectedRisksForNode() : ProcessFlowRisk[] {
            let allRisks = VueSetup.store.state.currentProcessFlowFullData.RiskKeys.map(
                ele => VueSetup.store.state.currentProcessFlowFullData.Risks[ele])

            let alreadySelected = new Set<number>(this.risksForNode.map(ele => ele.Id))
            return allRisks.filter(ele => !alreadySelected.has(ele.Id))
        },
        selectedRiskIds() : number[] {
            let riskIds = [] as number[]
            for (let risk of this.selectedRisks) {
                riskIds.push(risk.Id)
            }
            return riskIds
        },
        risksForDeletionConfirmation() : string[] {
            return this.selectedRisks.map(ele => ele.Name)
        }
    },
    methods : {
        generateRiskUrl(risk : ProcessFlowRisk) {
            return createRiskUrl(
                //@ts-ignore
                this.$root.orgGroupId,
                risk.Id)
        },
        saveNewRisk(risk : ProcessFlowRisk) {
            let currentNodeId = this.currentNode.Id

            VueSetup.store.commit('setRisk', risk)
            VueSetup.store.commit('addRisksToNode', {
                nodeId: currentNodeId,
                riskIds: [risk.Id]
            })

            this.showHideCreateNewRisk = false
        },
        cancelNewRisk() {
            this.showHideCreateNewRisk = false
        },
        addExistingRisk(selectedRisks : ProcessFlowRisk[]) {
            let riskIds = [] as number[]
            for (let r of selectedRisks) {
                riskIds.push(r.Id)
            }

            if (riskIds.length == 0) {
                this.showHideAddExistingRisk = false
                return
            }

            let currentNodeId = this.currentNode.Id
            addExistingRisk(<TAddExistingRiskInput>{
                csrf: getCurrentCSRF(),
                nodeId: currentNodeId,
                riskIds: riskIds
            }).then((resp : TAddExistingRiskOutput) => {
                VueSetup.store.commit('addRisksToNode', {
                    nodeId: currentNodeId,
                    riskIds: riskIds
                })
                this.showHideAddExistingRisk = false
            }).catch((err) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        cancelAddRisk() {
            this.showHideAddExistingRisk = false
        },
        toggleSelection() {
            if (this.hasSelected) {
                this.selectedRisks = []
            } else {
                this.selectedRisks = this.risksForNode
            }
        },
        deleteSelectedRisks(global : boolean) {
            let currentNodeId = this.currentNode.Id
            deleteRisk(<TDeleteRiskInput>{
                csrf: getCurrentCSRF(),
                nodeId: currentNodeId,
                riskIds: this.selectedRiskIds,
                global: global
            }).then((resp : TDeleteRiskOutput) => {
                VueSetup.store.dispatch('deleteBatchRisks', {
                    nodeId: currentNodeId,
                    riskIds: this.selectedRiskIds,
                    global: global
                })
                this.selectedRisks = []
                this.showHideDeleteRisk = false
            }).catch((err) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        editRisk(risk : ProcessFlowRisk) {
            this.editRiskData = {...risk}
            this.showHideEditRisk = true
            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editRisk.clearForm()
            })
        },
        cancelEditRisk() {
            this.showHideEditRisk = false
        },
        finishEditRisk(risk : ProcessFlowRisk) {
            VueSetup.store.commit('setRisk', risk)
            this.showHideEditRisk = false
        }
    },
    watch : {
        currentNode() {
            this.selectedRisks = []
        }
    }
})
</script>
