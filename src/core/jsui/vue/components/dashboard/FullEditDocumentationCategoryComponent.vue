<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Category: {{ currentCat.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`" >
                        {{ currentCat.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>

            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="4">
                        <create-new-control-documentation-category-form
                            edit-mode
                            :default-name="currentCat.Name"
                            :default-description="currentCat.Description"
                            :cat-id="currentCat.Id"
                            ref="editForm"
                            class="mb-4"
                        ></create-new-control-documentation-category-form>

                        <v-card class="mb-4">
                            <v-card-title>
                                <span class="mr-2">
                                    Input for Controls
                                </span>
                                <v-spacer></v-spacer>

                                <v-dialog persistent
                                          max-width="40%"
                                          v-model="showHideAddInput">
                                    <template v-slot:activator="{ on }">
                                        <v-btn color="primary" icon v-on="on">
                                            <v-icon>mdi-plus</v-icon>
                                        </v-btn>
                                    </template>
                                    
                                    <add-document-category-to-control-form
                                        :is-input="true"
                                        @do-cancel="showHideAddInput = false"
                                        @do-save="addInputControl"
                                        :fixed-cat="currentCat"
                                        :control-choices="allControls"
                                    ></add-document-category-to-control-form>
                                </v-dialog>

                            </v-card-title>
                            <v-divider></v-divider>

                            <control-table
                                :resources="inputControls"
                                use-crud-delete
                                @delete="deleteInputControl(currentCat, arguments[0])"
                            ></control-table>
                        </v-card>

                        <v-card>
                            <v-card-title>
                                <span class="mr-2">
                                    Output for Controls
                                </span>
                                <v-spacer></v-spacer>

                                <v-dialog persistent
                                          max-width="40%"
                                          v-model="showHideAddOutput">
                                    <template v-slot:activator="{ on }">
                                        <v-btn color="primary" icon v-on="on">
                                            <v-icon>mdi-plus</v-icon>
                                        </v-btn>
                                    </template>

                                    <add-document-category-to-control-form
                                        :is-input="false"
                                        @do-cancel="showHideAddOutput = false"
                                        @do-save="addOutputControl"
                                        :fixed-cat="currentCat"
                                        :control-choices="allControls"
                                    ></add-document-category-to-control-form>
                                </v-dialog>
                            </v-card-title>
                            <v-divider></v-divider>

                            <control-table
                                :resources="outputControls"
                                use-crud-delete
                                @delete="deleteOutputControl(currentCat, arguments[0])"
                            ></control-table>
                        </v-card>
                    </v-col>

                    <v-col cols="8">
                        <documentation-category-viewer :cat-id="currentCat.Id">
                        </documentation-category-viewer>
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { ControlDocumentationCategory } from '../../../ts/controls'
import { TGetDocCatOutput, getDocumentCategory } from '../../../ts/api/apiControlDocumentation'
import { TAllControlOutput, getAllControls} from '../../../ts/api/apiControls'
import { linkControlToDocumentCategory, unlinkControlFromDocumentCategory } from '../../../ts/api/apiControls'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import ControlTable from '../../generic/ControlTable.vue'
import DocumentationCategoryViewer from './DocumentationCategoryViewer.vue'
import CreateNewControlDocumentationCategoryForm from './CreateNewControlDocumentationCategoryForm.vue'
import AddDocumentCategoryToControlForm from '../../generic/AddDocumentCategoryToControlForm.vue'

@Component({
    components: {
        DocumentationCategoryViewer,
        CreateNewControlDocumentationCategoryForm,
        ControlTable,
        AddDocumentCategoryToControlForm
    }
})
export default class FullEditDocumentationCategoryComponent extends Vue {
    expandDescription : boolean = false
    currentCat : ControlDocumentationCategory | null = null
    inputControls: ProcessFlowControl[] = []
    outputControls: ProcessFlowControl[] = []
    allControls: ProcessFlowControl[] | null = null
    showHideAddInput: boolean = false
    showHideAddOutput: boolean = false

    get ready() : boolean {
        return !!this.currentCat && (this.allControls != null)
    }

    $refs!: {
        editForm : CreateNewControlDocumentationCategoryForm
    }
    
    mounted() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getDocumentCategory({
            orgId: PageParamsStore.state.organization!.Id,
            catId: resourceId,
        }).then((resp : TGetDocCatOutput) => {
            this.inputControls = resp.data.InputFor
            this.outputControls = resp.data.OutputFor
            this.currentCat = resp.data.Cat
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })

        getAllControls({
            orgName: PageParamsStore.state.organization!.OktaGroupName
        }).then((resp : TAllControlOutput) => {
            this.allControls = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    addIoControl(cat : ControlDocumentationCategory, ctrl : ProcessFlowControl, isInput: boolean) {
        linkControlToDocumentCategory({
            controlId: ctrl.Id,
            orgId: PageParamsStore.state.organization!.Id,
            catId: cat.Id,
            isInput: isInput
        }).then(() => {
            this.showHideAddInput = false
            this.showHideAddOutput = false
            if (isInput) {
                this.inputControls.push(ctrl)
            } else {
                this.outputControls.push(ctrl)
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
    }

    addInputControl(cat : ControlDocumentationCategory, ctrl : ProcessFlowControl) {
        this.addIoControl(cat, ctrl, true)
    }

    addOutputControl(cat : ControlDocumentationCategory, ctrl : ProcessFlowControl) {
        this.addIoControl(cat, ctrl, false)
    }

    deleteIoControl(cat : ControlDocumentationCategory, ctrl : ProcessFlowControl, isInput: boolean) {
        unlinkControlFromDocumentCategory({
            controlId: ctrl.Id,
            orgId: PageParamsStore.state.organization!.Id,
            catId: cat.Id,
            isInput: isInput
        }).then(() => {
            if (isInput) {
                this.inputControls = this.inputControls.filter((ele : ProcessFlowControl) => ele.Id != ctrl.Id)
            } else {
                this.outputControls = this.outputControls.filter((ele : ProcessFlowControl) => ele.Id != ctrl.Id)
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
    }

    deleteInputControl(cat : ControlDocumentationCategory, ctrl : ProcessFlowControl) {
        this.deleteIoControl(cat, ctrl, true)
    }

    deleteOutputControl(cat : ControlDocumentationCategory, ctrl : ProcessFlowControl) {
        this.deleteIoControl(cat, ctrl, false)
    }

    @Watch('ready')
    onReady() {
        Vue.nextTick(() => {
            this.$refs.editForm.clearForm()
        })
    }
}

</script>
