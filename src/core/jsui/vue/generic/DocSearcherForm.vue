<template>

<v-card>
    <v-card-title>
        Document Searcher
    </v-card-title>
    <v-divider></v-divider>

    <v-form @submit.prevent v-model="formValid" class="ma-4">
        <document-category-search-form-component
            v-model="chosenCat"
            load-cats
            :rules="[rules.required]"
        >
        </document-category-search-form-component>

        <div v-if="!!chosenCat">
            <doc-file-table
                v-model="selectedFiles"
                v-if="!loadingFiles"
                selectable
                multi
                :resources="currentCatFiles"
            >
            </doc-file-table>

            <v-row align="center" justify="center" class="py-4" v-else>
                <v-progress-circular indeterminate size="64"></v-progress-circular>
            </v-row>
        </div>
    </v-form>

    <v-card-actions>
        <v-btn color="error" @click="onCancel">
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn type="submit" color="primary" :disabled="!canSelect" @click="onSelect">
            Select
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DocumentCategorySearchFormComponent from './DocumentCategorySearchFormComponent.vue'
import DocFileTable from './DocFileTable.vue'
import { ControlDocumentationFile, ControlDocumentationCategory } from '../../ts/controls'
import { Watch } from 'vue-property-decorator'
import * as rules from '../../ts/formRules'
import { getControlDocuments, TGetControlDocumentsOutput } from '../../ts/api/apiControlDocumentation'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

const Props = Vue.extend({
    props: {
        excludeFiles: {
            type: Array,
            default: []
        }
    }
})

@Component({
    components: {
        DocumentCategorySearchFormComponent,
        DocFileTable
    }
})
export default class DocSearcher extends Props {
    chosenCat: ControlDocumentationCategory | null = null
    currentCatFiles : ControlDocumentationFile[] = []
    selectedFiles : ControlDocumentationFile[] = []
    loadingFiles: boolean = false
    formValid: boolean = false
    rules: any = rules

    get canSelect() : boolean {
        return this.formValid && this.selectedFiles.length > 0
    }

    @Watch('chosenCat')
    reloadFiles() {
        if (!this.chosenCat) {
            return
        }

        this.loadingFiles = true
        this.currentCatFiles = []
        getControlDocuments({
            catId: this.chosenCat!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetControlDocumentsOutput) => {
            let excludeSet = new Set<number>(this.excludeFiles.map((ele : any) => ele.Id))
            this.currentCatFiles = resp.data.Files.filter(
                (ele : ControlDocumentationFile) => !excludeSet.has(ele.Id)
            )
            this.loadingFiles = false
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

    onSelect() {
        this.$emit('do-select', this.selectedFiles)
        this.chosenCat = null
    }

    onCancel() {
        this.$emit('do-cancel')
        this.chosenCat = null
    }
}

</script>
