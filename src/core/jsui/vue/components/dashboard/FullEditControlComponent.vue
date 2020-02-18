<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Control: {{ fullControlData.Control.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ fullControlData.Control.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="8">
                        <create-new-control-form ref="editControl"
                                                 :node-id="-1"
                                                 :risk-id="-1"
                                                 :edit-mode="true"
                                                 :control="fullControlData.Control"
                                                 :staged-edits="true"
                                                 @do-save="onEditControl">
                        </create-new-control-form>
                    </v-col>

                    <v-col cols="4">
                        <v-card class="mb-4">
                            <v-card-title>
                                <span class="mr-2">
                                    Input Documentation
                                </span>
                                <v-spacer></v-spacer>

                                <v-dialog persistent
                                          max-width="40%"
                                          v-model="showHideAddInputDocCat">
                                    <template v-slot:activator="{ on }">
                                        <v-btn color="primary" icon v-on="on">
                                            <v-icon>mdi-plus</v-icon>
                                        </v-btn>
                                    </template>
                                    
                                    <add-document-category-to-control-form
                                        :is-input="true"
                                        @do-cancel="showHideAddInputDocCat = false"
                                        @do-save="addInputCat"
                                        :fixed-control="fullControlData.Control"
                                        :cat-choices="availableDocCats"
                                    ></add-document-category-to-control-form>
                                </v-dialog>
                            </v-card-title>
                            <v-divider></v-divider>

                            <documentation-table
                                :resources="fullControlData.InputDocCats"
                                use-crud-delete
                                @delete="deleteInputCat"
                            ></documentation-table>
                        </v-card>

                        <v-card class="mb-4">
                            <v-card-title>
                                <span class="mr-2">
                                    Output Documentation
                                </span>
                                <v-spacer></v-spacer>

                                <v-dialog persistent
                                          max-width="40%"
                                          v-model="showHideAddOutputDocCat">
                                    <template v-slot:activator="{ on }">
                                        <v-btn color="primary" icon v-on="on">
                                            <v-icon>mdi-plus</v-icon>
                                        </v-btn>
                                    </template>
                                    
                                    <add-document-category-to-control-form
                                        :is-input="true"
                                        @do-cancel="showHideAddOutputDocCat = false"
                                        @do-save="addOutputCat"
                                        :fixed-control="fullControlData.Control"
                                        :cat-choices="availableDocCats"
                                    ></add-document-category-to-control-form>
                                </v-dialog>

                            </v-card-title>
                            <v-divider></v-divider>

                            <documentation-table
                                :resources="fullControlData.OutputDocCats"
                                use-crud-delete
                                @delete="deleteOutputCat"
                            ></documentation-table>
                        </v-card>

                        <v-card class="mb-4">
                            <v-card-title>
                                Related Nodes
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-list two-line>
                                <v-list-item v-for="(item, index) in fullControlData.Nodes"
                                             :key="index"
                                >
                                    <v-list-item-content>
                                        <v-list-item-title>
                                            {{ item.Name }}
                                        </v-list-item-title>

                                        <v-list-item-subtitle>
                                            {{ item.Description }}
                                        </v-list-item-subtitle>
                                    </v-list-item-content>
                                </v-list-item>
                            </v-list>
                        </v-card>

                        <v-card>
                            <v-card-title>
                                Related Risks
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-list two-line>
                                <v-list-item v-for="(item, index) in fullControlData.Risks"
                                             :key="index"
                                             :href="generateRiskUrl(item.Id)"
                                >
                                    <v-list-item-content>
                                        <v-list-item-title>
                                            {{ item.Name }}
                                        </v-list-item-title>

                                        <v-list-item-subtitle>
                                            {{ item.Description }}
                                        </v-list-item-subtitle>
                                    </v-list-item-content>
                                </v-list-item>
                            </v-list>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import CreateNewControlForm from './CreateNewControlForm.vue'
import { FullControlData } from '../../../ts/controls'
import { getSingleControl, TSingleControlInput, TSingleControlOutput } from '../../../ts/api/apiControls'
import { createRiskUrl, contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import DocumentationTable from '../../generic/DocumentationTable.vue'
import { ControlDocumentationCategory } from '../../../ts/controls'
import { linkControlToDocumentCategory, unlinkControlFromDocumentCategory } from '../../../ts/api/apiControls'
import { getAllDocumentationCategories, TGetAllDocumentationCategoriesOutput } from '../../../ts/api/apiControlDocumentation'
import AddDocumentCategoryToControlForm from '../../generic/AddDocumentCategoryToControlForm.vue'

export default Vue.extend({
    data: () => ({
        expandDescription: false,
        fullControlData: null as FullControlData | null,
        showHideNewCat: false,
        showHideEditCat : false,
        showHideDeleteCat : false,
        showHideAddInputDocCat: false,
        showHideAddOutputDocCat: false,
        availableDocCats: null as ControlDocumentationCategory[] | null
    }),
    computed: {
        ready() : boolean {
            return (this.fullControlData != null && this.availableDocCats != null)
        }
    },
    watch: {
        ready() {
            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editControl.clearForm()
            })
        }
    },
    methods: {
        refreshData() {
            let data = window.location.pathname.split('/')
            let controlId = Number(data[data.length - 1])

            getSingleControl(<TSingleControlInput>{
                controlId: controlId,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TSingleControlOutput) => {
                this.fullControlData = resp.data
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })

            getAllDocumentationCategories({
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TGetAllDocumentationCategoriesOutput) => {
                this.availableDocCats = resp.data
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
        onEditControl(control : ProcessFlowControl) {
            this.fullControlData!.Control.Name = control.Name
            this.fullControlData!.Control.Description = control.Description
            this.fullControlData!.Control.ControlTypeId = control.ControlTypeId
            this.fullControlData!.Control.FrequencyType = control.FrequencyType
            this.fullControlData!.Control.FrequencyInterval = control.FrequencyInterval
            this.fullControlData!.Control.OwnerId = control.OwnerId
            this.fullControlData!.Control.Manual = control.Manual

            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editControl.clearForm()
            })
        },
        generateRiskUrl(riskId : number) : string {
            return createRiskUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                riskId)
        },

        addIoCat(cat : ControlDocumentationCategory, control : ProcessFlowControl, isInput: boolean) {
            linkControlToDocumentCategory({
                controlId: control.Id,
                orgId: PageParamsStore.state.organization!.Id,
                catId: cat.Id,
                isInput: isInput
            }).then(() => {
                this.showHideAddInputDocCat = false
                this.showHideAddOutputDocCat = false
                if (isInput) {
                    this.fullControlData!.InputDocCats.push(cat)
                } else {
                    this.fullControlData!.OutputDocCats.push(cat)
                }
            }).catch((err: any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        addInputCat(cat : ControlDocumentationCategory, control : ProcessFlowControl) {
            this.addIoCat(cat, control, true)
        },
        addOutputCat(cat : ControlDocumentationCategory, control : ProcessFlowControl) {
            this.addIoCat(cat, control, false)
        },

        deleteIoCat(cat : ControlDocumentationCategory, control : ProcessFlowControl, isInput: boolean) {
            unlinkControlFromDocumentCategory({
                controlId: control.Id,
                orgId: PageParamsStore.state.organization!.Id,
                catId: cat.Id,
                isInput: isInput
            }).then(() => {
                if (isInput) {
                    this.fullControlData!.InputDocCats = this.fullControlData!.InputDocCats.filter((ele : ControlDocumentationCategory) =>
                        ele.Id != cat.Id)
                } else {
                    this.fullControlData!.OutputDocCats = this.fullControlData!.OutputDocCats.filter((ele : ControlDocumentationCategory) =>
                        ele.Id != cat.Id)
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
        deleteInputCat(cat : ControlDocumentationCategory) {
            this.deleteIoCat(cat, this.fullControlData!.Control, true)
        },
        deleteOutputCat(cat : ControlDocumentationCategory) {
            this.deleteIoCat(cat, this.fullControlData!.Control, false)
        },
    },
    components: {
        CreateNewControlForm,
        DocumentationTable,
        AddDocumentCategoryToControlForm
    },
    mounted() {
        this.refreshData()
    }
})

</script>

<style scoped>

>>>.v-select__selections input {
    display: none !important;
}

>>>.v-tabs .v-slide-group {
    display: none !important;
}

</style>
