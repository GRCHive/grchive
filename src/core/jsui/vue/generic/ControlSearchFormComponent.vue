<template>
    <v-autocomplete
        :value="value"
        @input="onInput"
        label="Control"
        :items="items"
        :readonly="readonly"
        :rules="rules"
        filled
    >
    </v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { PageParamsStore } from '../../ts/pageParams'
import { getAllControls, TAllControlOutput } from '../../ts/api/apiControls'
import { NullControlFilterData } from '../../ts/controls'
import { contactUsUrl } from '../../ts/url'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => null as ProcessFlowControl | null
        },
        readonly: {
            type: Boolean,
            default: false
        },
        rules: {
            type: Array,
            default: []
        }
    }
})

@Component
export default class ControlSearchFormComponent extends Props {
    allControls : ProcessFlowControl[] = []

    get items() : any[] {
        return this.allControls.map((ele : ProcessFlowControl) => ({
            text: ele.Name,
            value: ele
        }))
    }


    compare(a : ProcessFlowControl | null, b : ProcessFlowControl | null) : boolean {
        if (!a || !b) {
            return false
        }

        return a.Id == b.Id
    }

    refreshControls() {
        getAllControls({
            orgName: PageParamsStore.state.organization!.OktaGroupName,
            filter: NullControlFilterData,
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
    
    onInput(c : ProcessFlowControl) {
        this.$emit('input', c)
    }

    mounted() {
        this.refreshControls()
    }
}

</script>
