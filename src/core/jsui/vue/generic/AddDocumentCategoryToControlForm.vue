<template>

<v-card>
    <v-card-title>
        Add {{ inputOutputText }} Document Category to Control
    </v-card-title>

    <v-form class="ma-4" v-model="formValid">
        <v-autocomplete
            label="Control"
            v-model="selectedControl"
            :rules="[rules.required]"
            :value-comparator="compareControls"
            :disabled="!!fixedControl"
            :items="availableControls"
        ></v-autocomplete>

        <v-autocomplete
            label="Document Category"
            v-model="selectedCat"
            :rules="[rules.required]"
            :value-comparator="compareDocumentationCategories"
            :disabled="!!fixedCat"
            :items="availableCategories"
        ></v-autocomplete>
    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
        >
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="success"
            @click="save"
            :disabled="!formValid"
        >
            Save
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ControlDocumentationCategory, compareDocumentationCategories, compareControls } from '../../ts/controls'
import * as rules from '../../ts/formRules'

const Props = Vue.extend({
    props: {
        fixedCat: {
            type: Object as () => ControlDocumentationCategory | null,
            default: null
        },
        fixedControl: {
            type: Object as () => ProcessFlowControl | null,
            default: null
        },
        catChoices: {
            type: Array,
            default: () => []
        },
        controlChoices: {
            type: Array,
            default: () => []
        },
        isInput: Boolean
    }
})

@Component
export default class AddDocumentCategoryToControlForm extends Props {
    formValid: boolean = false
    rules = rules

    selectedCat: ControlDocumentationCategory | null = null
    selectedControl: ProcessFlowControl | null = null

    get inputOutputText() : string {
        return this.isInput ? "Input" : "Output"
    }

    cancel() {
        this.$emit('do-cancel')
    }

    save() {
        this.$emit('do-save', this.selectedCat, this.selectedControl)
    }

    mounted() {
        this.selectedCat = this.fixedCat
        this.selectedControl = this.fixedControl
    }

    get availableCategories() : any[] {
        let mapFn = (ele : any) => ({
            text: ele.Name,
            value: ele
        })

        if (!!this.fixedCat) {
            return [mapFn(this.fixedCat)]
        }

        return this.catChoices.map(mapFn)
    }

    get availableControls() : any[] {
        let mapFn = (ele : any) => ({
            text: ele.Name,
            value: ele
        })

        if (!!this.fixedControl) {
            return [mapFn(this.fixedControl)]
        }
        return this.controlChoices.map(mapFn)
    }

    compareDocumentationCategories = compareDocumentationCategories
    compareControls = compareControls
}

</script>
