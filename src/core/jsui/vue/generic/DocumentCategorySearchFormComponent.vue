<template>
    <v-autocomplete
        label="Document Category"
        filled
        cache-items
        :items="displayItems"
        hide-no-data
        hide-selected
        :clearable="!readonly && !disabled"
        :value="value"
        @change="change"
        :disabled="disabled"
        :readonly="readonly"
        :value-comparator="compare"
        :rules="rules"
        :loading="!finishedLoading"
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ControlDocumentationCategory } from '../../ts/controls'
import { getAllDocumentationCategories, TGetAllDocumentationCategoriesOutput } from '../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'

const Props = Vue.extend({
    props: {
        value: {
            type: Object as () => ControlDocumentationCategory
        },
        availableCats: {
            type: Array,
        },
        initialCatId: {
            type: Number,
            default: -1
        },
        disabled: {
            type: Boolean,
            default: false
        },
        rules: {
            type: Array,
            default: () => []
        },
        loadCats: {
            type: Boolean,
            default: false
        },
        readonly: {
            type: Boolean,
            default: false
        }
    }
})

@Component
export default class DocumentCategorySearchFormComponent extends Props {
    loadedCategories : ControlDocumentationCategory[] = []
    finishedLoading: boolean = false

    extractId(a : ControlDocumentationCategory | null) : number {
        if (!a) {
            return -1
        }

        return a!.Id
    }

    compare(a : ControlDocumentationCategory | null, b : ControlDocumentationCategory | null) : boolean {
        if (!a || !b) {
            return false
        }

        let aId : number = this.extractId(a)
        let bId : number = this.extractId(b)
        return aId == bId
    }

    get displayItems() : Object[] {
        return this.displayCats.map((ele : any) => ({
            text: ele.Name,
            value: ele 
        }))
    }

    change(val : ControlDocumentationCategory) {
        this.$emit('input', val)
    }

    get displayCats() : ControlDocumentationCategory[] {
        if (this.loadCats) {
            return this.loadedCategories
        } else {
            return this.availableCats as ControlDocumentationCategory[]
        }
    }

    loadAllCategories()  {
        getAllDocumentationCategories({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetAllDocumentationCategoriesOutput) => {
            this.loadedCategories = resp.data
            this.finishedLoading = true

            if (this.initialCatId >= 0) {
                let idx = this.loadedCategories.findIndex((ele : ControlDocumentationCategory) => ele.Id == this.initialCatId)
                if (idx != -1) {
                    this.change(this.loadedCategories[idx])
                }
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

    mounted() {
        if (this.loadCats) {
            this.loadAllCategories()
        } else {
            this.finishedLoading = true
        }
    }
}


</script>
