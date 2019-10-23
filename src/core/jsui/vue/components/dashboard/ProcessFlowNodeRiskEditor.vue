<template>
    <section>
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
                    <delete-risk-form
                        :risks-to-delete="selectedRisks"
                        v-on:do-cancel="showHideDeleteRisk = false"
                        v-on:do-delete="deleteSelectedRisks">
                    </delete-risk-form>
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
                    <add-existing-risk-form
                        :preselected-risks="risksForNode"
                        @do-select="addExistingRisk"
                        @do-cancel="cancelAddRisk">
                    </add-existing-risk-form>
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
import AddExistingRiskForm from './AddExistingRiskForm.vue'
import DeleteRiskForm from './DeleteRiskForm.vue'
import { deleteRisk, addExistingRisk } from '../../../ts/api/apiRisks'
import { contactUsUrl } from '../../../ts/url'

export default Vue.extend({
    data : () => ({
        showHideDeleteRisk : false,
        showHideCreateNewRisk : false,
        showHideAddExistingRisk : false,
        selectedRisks : [] as ProcessFlowRisk[]
    }),
    components : {
        CreateNewRiskForm,
        DeleteRiskForm,
        AddExistingRiskForm
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
        selectedRiskIds() : number[] {
            let riskIds = [] as number[]
            for (let risk of this.selectedRisks) {
                riskIds.push(risk.Id)
            }
            return riskIds
        }
    },
    methods : {
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
                //@ts-ignore
                csrf: this.$root.csrf,
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
                //@ts-ignore
                csrf: this.$root.csrf,
                nodeId: currentNodeId,
                riskIds: this.selectedRiskIds,
                global: global
            }).then((resp : TDeleteRiskOutput) => {
                VueSetup.store.dispatch('deleteBatchRisks', {
                    nodeId: currentNodeId,
                    riskIds: this.selectedRiskIds,
                    global: global
                })
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
        }
    },
    watch : {
        currentNode() {
            this.selectedRisks = []
        }
    }
})
</script>
