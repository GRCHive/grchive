<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Documentation
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
                <v-dialog v-model="showHideNew"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <create-new-control-documentation-category-form
                        @do-cancel="showHideNew = false"
                        @do-save="saveNewControlDocCategory">
                    </create-new-control-documentation-category-form>
                </v-dialog>
            </v-list-item-action>

        </v-list-item>

        <v-divider></v-divider>

        <documentation-table
            :resources="docCats"
            :search="filterText"
            use-crud-delete
            confirm-delete
            @delete="deleteControlDocCategory"
        ></documentation-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DocumentationTable from '../../generic/DocumentationTable.vue'
import { ControlDocumentationCategory } from '../../../ts/controls'
import CreateNewControlDocumentationCategoryForm from './CreateNewControlDocumentationCategoryForm.vue'
import { getAllDocumentationCategories, TGetAllDocumentationCategoriesOutput } from '../../../ts/api/apiControlDocumentation'
import { deleteControlDocCat } from '../../../ts/api/apiControlDocumentation'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

@Component({
    components: {
        DocumentationTable,
        CreateNewControlDocumentationCategoryForm
    }
})
export default class DashboardDocumentationList extends Vue {
    filterText: string = ""
    docCats : ControlDocumentationCategory[] = []
    showHideNew: boolean = false

    saveNewControlDocCategory(cat : ControlDocumentationCategory) {
        this.showHideNew = false
        this.docCats.unshift(cat)
    }

    deleteControlDocCategory(cat : ControlDocumentationCategory) {
        deleteControlDocCat({
            catId: cat.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then(() => {
            this.docCats = this.docCats.filter((ele : ControlDocumentationCategory) => ele.Id != cat.Id)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        getAllDocumentationCategories({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TGetAllDocumentationCategoriesOutput) => {
            this.docCats = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }
}

</script>
