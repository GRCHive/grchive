<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <v-dialog v-model="showHideEditCat" persistent max-width="40%">
            <create-new-control-documentation-category-form
                ref="editControlDocCat"
                :control="fullControlData.Control"
                :edit-mode="true"
                :cat-id="currentDocumentCategory.Id"
                :default-name="currentDocumentCategory.Name"
                :default-description="currentDocumentCategory.Description"
                @do-cancel="cancelEditControlDocCategory"
                @do-save="saveEditControlDocCategory">
            </create-new-control-documentation-category-form>
        </v-dialog>

        <v-dialog v-model="showHideDeleteCat" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="documentation categories"
                :items-to-delete="[currentDocumentCategory.Name]"
                v-on:do-cancel="showHideDeleteCat = false"
                v-on:do-delete="deleteSelectedCategories"
                :use-global-deletion="false"
                :force-global-deletion="true">
            </generic-delete-confirmation-form>
        </v-dialog>

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
                                    Documentation
                                </span>
                                <v-spacer></v-spacer>
                                <v-select
                                    max-width="50%"
                                    :hide-details="true"
                                    :items="categoriesForSelect"
                                    v-model="currentDocumentCategory"
                                    class="pa-0 ma-0">

                                    <template v-slot:prepend>
                                        <v-btn small color="error" @click="doDeleteControlDocCat" :disabled="!hasCategories">
                                            Delete
                                        </v-btn>
                                    </template>
                                    <template v-slot:append-outer>
                                        <v-btn small color="warning" @click="doEditControlDocCat" class="mr-2" :disabled="!hasCategories">
                                            Edit
                                        </v-btn>

                                        <v-dialog v-model="showHideNewCat" persistent max-width="40%">
                                            <template v-slot:activator="{ on }">
                                                <v-btn small color="primary" v-on="on">
                                                    New
                                                </v-btn>
                                            </template>
                                            <create-new-control-documentation-category-form
                                                :control="fullControlData.Control"
                                                @do-cancel="cancelNewControlDocCategory"
                                                @do-save="saveNewControlDocCategory">
                                            </create-new-control-documentation-category-form>
                                        </v-dialog>
                                    </template>
                                </v-select>
                            </v-card-title>
                            <v-divider></v-divider>
                            <v-tabs :value="currentDocumentCategory">
                                <v-tab-item
                                    v-for="(item, index) in fullControlData.DocumentCategories"
                                    :key="index"
                                    :value="item"
                                >
                                    <documentation-category-viewer
                                        :cat-id="item.Id">
                                    </documentation-category-viewer>
                                </v-tab-item>
                            </v-tabs>
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
import CreateNewControlDocumentationCategoryForm from './CreateNewControlDocumentationCategoryForm.vue'
import DocumentationCategoryViewer from './DocumentationCategoryViewer.vue'
import { FullControlData } from '../../../ts/controls'
import { getSingleControl, TSingleControlInput, TSingleControlOutput } from '../../../ts/api/apiControls'
import { createRiskUrl, contactUsUrl } from '../../../ts/url'
import { ControlDocumentationCategory } from '../../../ts/controls'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { deleteControlDocCat, TDeleteControlDocCatInput, TDeleteControlDocCatOutput } from '../../../ts/api/apiControlDocumentation'

export default Vue.extend({
    data: () => ({
        expandDescription: false,
        ready: false,
        fullControlData: Object() as FullControlData,
        showHideNewCat: false,
        showHideEditCat : false,
        showHideDeleteCat : false,
        currentDocumentCategory: Object() as ControlDocumentationCategory
    }),
    computed: {
        categoriesForSelect() : Object[] {
            return this.fullControlData.DocumentCategories.map((ele) => ({
                text: ele.Name,
                value: ele
            }))
        },
        hasCategories() : boolean {
            return this.fullControlData.DocumentCategories.length > 0
        }
    },
    methods: {
        resetSelectedCategory() {
            if (this.fullControlData.DocumentCategories.length > 0) {
                this.currentDocumentCategory = this.fullControlData.DocumentCategories[0]
            } else {
                this.currentDocumentCategory = Object() as ControlDocumentationCategory
            }
        },
        refreshData() {
            let data = window.location.pathname.split('/')
            let controlId = Number(data[data.length - 1])

            getSingleControl(<TSingleControlInput>{
                //@ts-ignore
                csrf: this.$root.csrf,
                controlId: controlId
            }).then((resp : TSingleControlOutput) => {
                this.fullControlData = resp.data
                this.ready = true
                this.resetSelectedCategory()

                Vue.nextTick(() => {
                    //@ts-ignore
                    this.$refs.editControl.clearForm()
                })
            }).catch((err : any) => {
                window.location.replace('/404')
            })
        },
        onEditControl(control : ProcessFlowControl) {
            this.fullControlData.Control.Name = control.Name
            this.fullControlData.Control.Description = control.Description
            this.fullControlData.Control.ControlTypeId = control.ControlTypeId
            this.fullControlData.Control.FrequencyType = control.FrequencyType
            this.fullControlData.Control.FrequencyInterval = control.FrequencyInterval
            this.fullControlData.Control.OwnerId = control.OwnerId

            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editControl.clearForm()
            })
        },
        generateRiskUrl(riskId : number) : string {
            return createRiskUrl(
                //@ts-ignore
                this.$root.orgGroupId,
                riskId)
        },
        saveNewControlDocCategory(cat : ControlDocumentationCategory) {
            this.showHideNewCat = false
            this.fullControlData.DocumentCategories.push(cat)
        },
        cancelNewControlDocCategory() {
            this.showHideNewCat = false
        },
        doEditControlDocCat() {
            this.showHideEditCat = true
            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editControlDocCat.clearForm()
            })
        },
        cancelEditControlDocCategory() {
            this.showHideEditCat = false
        },
        saveEditControlDocCategory(cat : ControlDocumentationCategory) {
            this.showHideEditCat = false
            this.currentDocumentCategory.Name = cat.Name
            this.currentDocumentCategory.Description = cat.Description
        },
        doDeleteControlDocCat() {
            this.showHideDeleteCat = true
        },
        deleteSelectedCategories() {
            deleteControlDocCat(<TDeleteControlDocCatInput>{
                //@ts-ignore
                csrf: this.$root.csrf,
                catId: this.currentDocumentCategory.Id,
            }).then(() => {
                this.showHideDeleteCat = false
                this.fullControlData.DocumentCategories.splice(
                    this.fullControlData.DocumentCategories.findIndex((ele) => 
                        ele.Name == this.currentDocumentCategory.Name),
                    1)
                this.resetSelectedCategory()
            }).catch(() => {
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
        CreateNewControlForm,
        CreateNewControlDocumentationCategoryForm,
        DocumentationCategoryViewer,
        GenericDeleteConfirmationForm,
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
