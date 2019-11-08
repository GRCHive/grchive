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
                            Create New
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
        <v-card
            v-for="(item, index) in filteredRisks"
            :key="index"
            class="my-2"
        >
            <v-list-item two-line @click="goToRisk(item.Id)">
                <v-list-item-content>
                    <v-list-item-title v-html="highlightText(item.Name)">
                    </v-list-item-title>
                    <v-list-item-subtitle v-html="highlightText(item.Description)">
                    </v-list-item-subtitle>
                </v-list-item-content>
                <v-spacer></v-spacer>
                <v-list-item-action>
                    <v-btn icon @click.stop="doDeleteRisk(item)" @mousedown.stop @mouseup.stop>
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </v-card>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { getAllRisks, TAllRiskInput, TAllRiskOutput } from '../../../ts/api/apiRisks'
import { deleteRisk, TDeleteRiskInput, TDeleteRiskOutput } from '../../../ts/api/apiRisks'
import { contactUsUrl, createRiskUrl } from '../../../ts/url'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'
import CreateNewRiskForm from './CreateNewRiskForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    data : () => ({
        allRisks: [] as ProcessFlowRisk[],
        filterText : "",
        showHideCreateNewRisk: false,
        showHideDeleteRisk: false,
        currentDeleteRisk : Object() as ProcessFlowRisk
    }),
    components: {
        CreateNewRiskForm,
        GenericDeleteConfirmationForm
    },
    computed: {
        filter() : (a : ProcessFlowRisk) => boolean {
            const filterText = this.filterText.trim()
            return (ele : ProcessFlowRisk) : boolean => {
                return ele.Name.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) ||
                    ele.Description.toLocaleLowerCase().includes(filterText.toLocaleLowerCase())
            }
        },
        filteredRisks() : ProcessFlowRisk[] {
            return this.allRisks.filter(this.filter)
        },
        currentRisksToDelete() : string[] {
            if (!this.showHideDeleteRisk) {
                return []
            }
            return [this.currentDeleteRisk.Name]
        }
    },
    methods: {
        highlightText(input : string) : string {
            const safeInput = sanitizeTextForHTML(input)
            const useFilter = this.filterText.trim()
            if (useFilter.length == 0) {
                return safeInput
            }
            return replaceWithMark(
                safeInput,
                sanitizeTextForHTML(useFilter))
        },
        generateRiskUrl(riskId : number) : string {
            return createRiskUrl(PageParamsStore.state.organization!.OktaGroupName, riskId)
        },
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
            this.allRisks.push(risk)
            this.showHideCreateNewRisk = false
        },
        cancelNewRisk() {
            this.showHideCreateNewRisk = false
        },
        goToRisk(riskId : number) {
            window.location.assign(this.generateRiskUrl(riskId))
        },
        doDeleteRisk(risk : ProcessFlowRisk) {
            this.currentDeleteRisk = risk
            this.showHideDeleteRisk = true
        },
        deleteSelectedRisks(global : boolean, items : string[]) {
            // assumption: global is true, items has length 1
            const idx = this.allRisks.findIndex(
                (ele) => items.includes(ele.Name))
            if (idx == -1) {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
                return // ???
            }

            const risk = this.allRisks[idx]

            deleteRisk(<TDeleteRiskInput>{
                nodeId: -1,
                riskIds: [risk.Id],
                global: true
            }).then((resp : TDeleteRiskOutput) => {
                this.allRisks.splice(idx, 1)
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
    mounted() {
        this.refreshRisks()
    }
})
</script>
