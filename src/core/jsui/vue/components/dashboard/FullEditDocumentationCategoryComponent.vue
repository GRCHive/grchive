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
                        ></create-new-control-documentation-category-form>
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
import { ControlDocumentationCategory } from '../../../ts/controls'
import { TGetDocCatOutput, getDocumentCategory } from '../../../ts/api/apiControlDocumentation'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import DocumentationCategoryViewer from './DocumentationCategoryViewer.vue'
import CreateNewControlDocumentationCategoryForm from './CreateNewControlDocumentationCategoryForm.vue'

@Component({
    components: {
        DocumentationCategoryViewer,
        CreateNewControlDocumentationCategoryForm
    }
})
export default class FullEditDocumentationCategoryComponent extends Vue {
    ready: boolean = false
    expandDescription : boolean = false
    currentCat : ControlDocumentationCategory | null = null

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
            this.currentCat = resp.data
            this.ready = true
            Vue.nextTick(() => {
                this.$refs.editForm.clearForm()
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
    }
}

</script>
