<template>
    <div class="ma-4">
        <v-dialog v-model="showHideDeleteFlow" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="process flows"
                :items-to-delete="currentFlowsToDelete"
                v-on:do-cancel="showHideDeleteFlow = false"
                v-on:do-delete="deleteFlow"
                :use-global-deletion="false"
                :force-global-deletion="true">
            </generic-delete-confirmation-form>
        </v-dialog>

        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Process Flows
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
                <v-dialog v-model="showHideCreateNewFlow" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            Create New
                        </v-btn>
                    </template>
                    <create-new-process-flow-form
                        @do-save="onCreateNewFlow"
                        @do-cancel="showHideCreateNewFlow = false"
                    >
                    </create-new-process-flow-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <v-list-item class="headerItem">
            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Process Flow
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Created Date
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Last Updated Date
                </v-list-item-title>
            </v-list-item-content>
            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn icon disabled></v-btn>
            </v-list-item-action>
        </v-list-item>

        <v-card
            v-for="(item, index) in filteredFlows"
            :key="index"
            class="my-2"
        >
            <v-list-item two-line @click.stop="goToFlow(item.Id)">
                <v-list-item-content>
                    <v-list-item-title v-html="highlightText(item.Name)">
                    </v-list-item-title>
                    <v-list-item-subtitle v-html="highlightText(item.Description)">
                    </v-list-item-subtitle>
                </v-list-item-content>
                <v-list-item-content>
                    <v-list-item-title>
                        {{ item.CreationTime.toDateString() }}
                    </v-list-item-title>
                </v-list-item-content>
                <v-list-item-content>
                    <v-list-item-title>
                        {{ item.LastUpdatedTime.toDateString() }}
                    </v-list-item-title>
                </v-list-item-content>

                <v-spacer></v-spacer>
                <v-list-item-action>
                    <v-btn icon @click.stop="startDeleteFlow(item)" @mousedown.stop @mouseup.stop>
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </v-card>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import CreateNewProcessFlowForm from './CreateNewProcessFlowForm.vue'
import { getAllProcessFlow, TGetAllProcessFlowInput, TGetAllProcessFlowOutput } from '../../../ts/api/apiProcessFlow'
import { deleteProcessFlow, TDeleteProcessFlowInput, TDeleteProcessFlowOutput } from '../../../ts/api/apiProcessFlow'
import { contactUsUrl, createFlowUrl } from '../../../ts/url'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    data : () => ({
        showHideCreateNewFlow: false,
        allFlows: [] as ProcessFlowBasicData[],
        filterText: "",
        showHideDeleteFlow: false,
        flowToDelete: Object() as ProcessFlowBasicData
    }),
    components: {
        CreateNewProcessFlowForm,
        GenericDeleteConfirmationForm
    },
    computed: {
        filter() : (a : ProcessFlowBasicData) => boolean {
            const filterText = this.filterText.trim()
            return (ele : ProcessFlowBasicData) : boolean => {
                return ele.Name.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) ||
                    ele.Description.toLocaleLowerCase().includes(filterText.toLocaleLowerCase())
            }
        },
        filteredFlows() : ProcessFlowBasicData[] {
            return this.allFlows.filter(this.filter)
        },
        currentFlowsToDelete() : string[] {
            if (!this.showHideDeleteFlow) {
                return []
            }
            return [this.flowToDelete.Name]
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
        onCreateNewFlow(data : ProcessFlowBasicData) {
            this.showHideCreateNewFlow = false
            this.allFlows.push(data)
        },
        refreshFlows() {
            getAllProcessFlow(<TGetAllProcessFlowInput>{
                requested: -1,
                organization: PageParamsStore.state.organization!.OktaGroupName
            }).then((resp : TGetAllProcessFlowOutput) => {
                this.allFlows = resp.data.Flows
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        goToFlow(flowId: number) {
            window.location.assign(createFlowUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                flowId))
        },
        startDeleteFlow(flow : ProcessFlowBasicData) {
            this.showHideDeleteFlow = true
            this.flowToDelete = flow
        },
        deleteFlow() {
            let processFlow = this.flowToDelete

            deleteProcessFlow(<TDeleteProcessFlowInput>{
                flowId: processFlow.Id,
            }).then((resp : TDeleteProcessFlowOutput) => {
                this.allFlows.splice(
                    this.allFlows.findIndex((ele : ProcessFlowBasicData) =>
                        ele.Id == processFlow.Id),
                    1)
                this.showHideDeleteFlow = false
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
    },
    mounted() {
        this.refreshFlows()
    }
})

</script>

<style scoped>

.headerItem {
    max-height: 30px !important;
    min-height: 30px !important;
}

</style>
