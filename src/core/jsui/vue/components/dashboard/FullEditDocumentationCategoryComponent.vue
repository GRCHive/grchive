<template>
    <div :class="contentOnly ? '' : 'ma-4'">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0" v-if="!contentOnly">
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

            <v-divider v-if="!contentOnly"></v-divider>

            <v-container fluid :class="contentOnly ? 'pa-0' : ''">
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="4">
                                <create-new-control-documentation-category-form
                                    edit-mode
                                    :default-name="currentCat.Name"
                                    :default-description="currentCat.Description"
                                    :cat-id="currentCat.Id"
                                    class="mb-4"
                                ></create-new-control-documentation-category-form>

                                <v-card class="mb-4">
                                    <v-card-title>
                                        <span class="mr-2">
                                            Related Resources
                                        </span>
                                        <v-spacer></v-spacer>
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <v-tabs>
                                        <v-tab v-if="!!relevantControls">Controls</v-tab>
                                        <v-tab-item v-if="!!relevantControls">
                                            <control-table
                                                :resources="relevantControls"
                                            ></control-table>
                                        </v-tab-item>
                                    </v-tabs>

                                </v-card>
                            </v-col>

                            <v-col cols="8">
                                <documentation-category-viewer :cat-id="currentCat.Id">
                                </documentation-category-viewer>
                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Audit Trail</v-tab>
                    <v-tab-item>
                        <audit-trail-viewer
                            resource-type="process_flow_control_documentation_categories"
                            :resource-id="`${currentCat.Id}`"
                            no-header
                        >
                        </audit-trail-viewer>
                    </v-tab-item>
                </v-tabs>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ControlDocumentationCategory, NullControlFilterData } from '../../../ts/controls'
import { TGetDocCatOutput, getDocumentCategory } from '../../../ts/api/apiControlDocumentation'
import { allControlDocCatLink, TAllControlDocCatLinkOutput } from '../../../ts/api/apiControlDocCatLinks'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import ControlTable from '../../generic/ControlTable.vue'
import DocumentationCategoryViewer from './DocumentationCategoryViewer.vue'
import CreateNewControlDocumentationCategoryForm from './CreateNewControlDocumentationCategoryForm.vue'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'

const Props = Vue.extend({
    props: {
        contentOnly : {
            type: Boolean,
            default: false
        },
        resourceId: {
            type: Number,
            default: -1
        }
    }
})

@Component({
    components: {
        DocumentationCategoryViewer,
        CreateNewControlDocumentationCategoryForm,
        ControlTable,
        AuditTrailViewer
    }
})
export default class FullEditDocumentationCategoryComponent extends Props {
    expandDescription : boolean = false
    currentCat : ControlDocumentationCategory | null = null
    showHideAddInput: boolean = false
    showHideAddOutput: boolean = false

    relevantControls : ProcessFlowControl[] | null = null

    get ready() : boolean {
        return !!this.currentCat
    }

    refreshRelevantControls() {
        allControlDocCatLink({
            catId: this.currentCat!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlDocCatLinkOutput) => {
            this.relevantControls = resp.data.Controls!
        })
    }
    
    mounted() {
        let resourceId = 0
        if (this.resourceId != -1) {
            resourceId = this.resourceId
        } else {
            let data = window.location.pathname.split('/')
            resourceId = Number(data[data.length - 1])
        }

        getDocumentCategory({
            orgId: PageParamsStore.state.organization!.Id,
            catId: resourceId,
            lean: false,
        }).then((resp : TGetDocCatOutput) => {
            this.currentCat = resp.data.Cat
            this.refreshRelevantControls()
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
}

</script>
