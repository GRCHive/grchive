<template>
    <section>
        <v-dialog v-model="showHideEditControl" persistent max-width="40%">
            <create-new-control-form
                ref="editControl"
                :edit-mode="true"
                :node-id="currentNode.Id"
                :risk-id="currentEditRiskId"
                :control="currentEditControl"
                @do-save="finishEditControl"
                @do-cancel="cancelEditControl">
            </create-new-control-form>
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
                <v-dialog v-model="showHideDeleteControl" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="error" :disabled="!hasSelected" v-on="on">
                            Delete
                            <v-icon small>mdi-delete</v-icon>
                        </v-btn>
                    </template>

                    <generic-delete-confirmation-form
                        item-name="controls"
                        :items-to-delete="controlsForDeletionConfirmation"
                        v-on:do-cancel="showHideDeleteControl = false"
                        v-on:do-delete="deleteSelectedControls"
                        :use-global-deletion="true">
                    </generic-delete-confirmation-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

        <v-dialog v-model="showHideNewControl" persistent max-width="40%">
            <create-new-control-form
                :node-id="currentNode.Id"
                :risk-id="currentRelevantRiskId"
                @do-cancel="onCancelNewControl"
                @do-save="onSaveNewControl">
            </create-new-control-form>
        </v-dialog>
        <v-dialog v-model="showHideExistingControl" persistent max-width="40%">
            <generic-add-existing-item-form
                item-name="Controls"
                :selectable-items="unselectedControlsForRisk(currentRelevantRiskId)"
                @do-select="onAddExistingControls"
                @do-cancel="onCancelExistingControls">
            </generic-add-existing-item-form>
        </v-dialog>

        <v-list two-line>
            <v-list-item-group multiple v-model="selectedControls">
                <v-list-group v-for="(item, index) in risksForNode"
                              :key="index"
                              class="pa-1"
                              :value="true">
                    <template v-slot:activator>
                        <v-list-item-action>
                            <v-menu offset-y>
                                <template v-slot:activator="{ on }">
                                    <v-btn icon v-on="on" @mousedown.stop @click.stop>
                                        <v-icon>mdi-plus</v-icon>
                                    </v-btn>
                                </template>
                                <v-list dense>
                                    <v-list-item @click="showNewControlDialog(item.Id)">
                                        New Control
                                    </v-list-item>
                                    <v-list-item @click="showExistingControlDialog(item.Id)">
                                        Existing Control
                                    </v-list-item>
                                </v-list>
                            </v-menu>
                        </v-list-item-action>
                        
                        <v-list-item-content>
                            <v-list-item-title>
                                Risk:&nbsp;{{ item.Name }}
                            </v-list-item-title>

                            <v-list-item-subtitle>
                                {{ item.Description }}
                            </v-list-item-subtitle>
                        </v-list-item-content>
                    </template>

                    <v-list-item two-line
                                 v-for="controlItem in controlsForRiskNode(item.Id)"
                                 :key="controlItem.control.Id"
                                 :value="controlItem">
                        <template v-slot:default="{active, toggle}">
                            <v-list-item-action>
                                <v-checkbox :input-value="active"
                                            :true-value="controlItem"
                                            @click="toggle">
                                </v-checkbox>
                            </v-list-item-action>
                            <v-list-item-content>
                                <v-list-item-title>
                                    {{ controlItem.control.Name }}
                                </v-list-item-title>

                                <v-list-item-subtitle>
                                    {{ controlItem.control.Description }}
                                </v-list-item-subtitle>
                            </v-list-item-content>

                            <v-list-item-action>
                                <v-btn icon @click.stop="editControl(controlItem)" @mousedown.stop>
                                    <v-icon>mdi-pencil</v-icon>
                                </v-btn>
                            </v-list-item-action>

                            <v-list-item-action>
                                <v-btn icon @click.stop @mousedown.stop :href="generateControlUrl(controlItem.control)" target="_blank">
                                    <v-icon>mdi-open-in-new</v-icon>
                                </v-btn>
                            </v-list-item-action>

                        </template>

                    </v-list-item>
                </v-list-group>
            </v-list-item-group>
        </v-list>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from "../../../ts/vueSetup"
import CreateNewControlForm from './CreateNewControlForm.vue'
import GenericAddExistingItemForm from './GenericAddExistingItemForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import {TDeleteControlInput, TDeleteControlOutput, deleteControls} from '../../../ts/api/apiControls'
import {TExistingControlInput, TExistingControlOutput, addExistingControls} from '../../../ts/api/apiControls'
import { contactUsUrl, createControlUrl } from '../../../ts/url'

export default Vue.extend({
    data : () => ({
        showHideDeleteControl : false,
        showHideNewControl : false,
        showHideExistingControl : false,
        showHideEditControl: false,
        currentRelevantRiskId : -1,
        selectedControls : [] as RiskControl[],
        currentEditRiskId: -1,
        currentEditControl : Object() as ProcessFlowControl
    }),
    components : {
        CreateNewControlForm,
        GenericAddExistingItemForm,
        GenericDeleteConfirmationForm
    },
    computed: {
        currentNode() : ProcessFlowNode {
            return VueSetup.store.getters.currentNodeInfo
        },
        risksForNode() : ProcessFlowRisk[] {
            return VueSetup.store.getters.risksForNode(this.currentNode.Id)
        },
        controlsForRiskNode() : (riskId : number) => RiskControl[] {
            let risks = this.risksForNode

            let cache = Object() as Record<number, RiskControl[]>
            for (let r of risks) {
                cache[r.Id] = VueSetup.store.getters.controlsForRiskNode(r.Id, this.currentNode.Id)
            }

            return (riskId: number) : RiskControl[] => {
                return cache[riskId]
            }
        },
        hasSelected() : boolean {
            return this.filteredRiskControl.length > 0
        },
        filteredRiskControl() : RiskControl[] {
            return this.selectedControls.filter(ele => typeof ele == 'object')
        },
        controlsForDeletionConfirmation() : string[] {
            return this.filteredRiskControl.map(ele => `${ele.control.Name} (${ele.risk.Name})`)
        },
        unselectedControlsForRisk() : (riskId : number) => ProcessFlowControl[] {
            let allControls = VueSetup.store.state.currentProcessFlowFullData.ControlKeys.map(
                ele => VueSetup.store.state.currentProcessFlowFullData.Controls[ele])

            let existSetCache = Object() as Record<number, Set<number>>
            for (let r of this.risksForNode) {
                existSetCache[r.Id] = new Set<number>(this.controlsForRiskNode(r.Id).map(
                    ele => ele.control.Id
                ))
            }

            return (riskId : number) : ProcessFlowControl[] => {
                if (riskId == -1) {
                    return []
                }
                return allControls.filter(ele => !existSetCache[riskId].has(ele.Id))
            }
        }
    },
    methods : {
        generateControlUrl(control : ProcessFlowControl) {
            return createControlUrl(
                //@ts-ignore
                this.$root.orgGroupId,
                control.Id)
        },
        showNewControlDialog(riskId : number) {
            this.currentRelevantRiskId = riskId
            this.showHideNewControl = true
        },
        showExistingControlDialog(riskId : number) {
            this.currentRelevantRiskId = riskId
            this.showHideExistingControl = true
        },
        onCancelNewControl() {
            this.showHideNewControl = false
        },
        onSaveNewControl(control : ProcessFlowControl, riskId : number) {
            VueSetup.store.commit('setControl', {control})
            VueSetup.store.commit('addControlToNode', {
                controlId: control.Id,
                nodeId: this.currentNode.Id
            })
            VueSetup.store.commit('addControlToRisk', {
                controlId: control.Id,
                riskId: riskId
            })
            this.showHideNewControl = false
        },
        onCancelExistingControls() {
            this.showHideExistingControl = false
        },
        onAddExistingControls(selectedControls: ProcessFlowControl[]) {
            let currentNodeId = this.currentNode.Id
            let currentRiskId = this.currentRelevantRiskId

            addExistingControls(<TExistingControlInput>{
                nodeId: currentNodeId,
                riskId: currentRiskId,
                controlIds: selectedControls.map(ele => ele.Id)
            }).then((resp : TExistingControlOutput) => {
                this.showHideExistingControl = false
                for (let control of selectedControls) {
                    VueSetup.store.commit('addControlToNode', {
                        controlId: control.Id,
                        nodeId: currentNodeId
                    })

                    VueSetup.store.commit('addControlToRisk', {
                        controlId: control.Id,
                        riskId: currentRiskId
                    })
                }
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
        toggleSelection() {
            if (this.hasSelected) {
                this.selectedControls = []
            } else {
                let set = new Set<RiskControl>() 
                for (let risk of this.risksForNode) {
                    this.controlsForRiskNode(risk.Id).forEach(ele => set.add(ele))
                }
                this.selectedControls = Array.from(set)
            }
        },
        editControl(control : RiskControl) {
            this.currentEditRiskId = control.risk.Id
            this.currentEditControl = control.control
            this.showHideEditControl = true
            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editControl.clearForm()
            })
        },
        deleteSelectedControls(global : boolean) {
            let currentNodeId = this.currentNode.Id
            let riskIds = this.filteredRiskControl.map(ele => ele.risk.Id)
            let controlIds = this.filteredRiskControl.map(ele => ele.control.Id)
            deleteControls(<TDeleteControlInput>{
                nodeId: currentNodeId,
                riskIds: riskIds,
                controlIds: controlIds,
                global: global
            }).then((resp : TDeleteControlOutput) => {
                this.selectedControls = []
                this.showHideDeleteControl = false
                VueSetup.store.dispatch('deleteBatchControls', {
                    nodeId: currentNodeId,
                    controlIds: controlIds,
                    riskIds: riskIds,
                    global: global
                })
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
        finishEditControl(control : ProcessFlowControl, riskId : number) {
            VueSetup.store.commit('setControl', {control})
            this.showHideEditControl = false
        },
        cancelEditControl() {
            this.showHideEditControl = false
        }
    },
    watch : {
        currentNode() {
            this.selectedControls = []
        },
    },
})

</script>
