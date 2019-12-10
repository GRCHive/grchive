<template>
    <v-autocomplete
        label="Document Category"
        filled
        cache-items
        :items="displayItems"
        hide-no-data
        hide-selected
        clearable
        :value="value"
        @change="change"
        :disabled="disabled"
        :value-comparator="compare"
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ControlDocumentationCategory } from '../../ts/controls'

const Props = Vue.extend({
    props: {
        value: {
            type: Object as () => ControlDocumentationCategory | number
        },
        availableCats: {
            type: Array,
        },
        disabled: {
            type: Boolean,
            default: false
        },
        idMode: {
            type: Boolean,
            default: false
        }
    }
})

@Component
export default class DocumentCategorySearchFormComponent extends Props {
    extractId(a : ControlDocumentationCategory | number) : number {
        if ((a as ControlDocumentationCategory).Id) {
            return (a as ControlDocumentationCategory).Id
        } else {
            return (a as number)
        }
    }

    compare(a : ControlDocumentationCategory | number, b : ControlDocumentationCategory | number) : boolean {
        let aId : number = this.extractId(a)
        let bId : number = this.extractId(b)
        return aId == bId
    }

    get displayItems() : Object[] {
        return this.availableCats.map((ele : any) => ({
            text: ele.Name,
            value: ele 
        }))
    }

    change(val : ControlDocumentationCategory) {
        if (this.idMode) {
            this.$emit('input', val.Id)
        } else {
            this.$emit('input', val)
        }
    }
}


</script>
