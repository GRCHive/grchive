<template>
    <v-select
        dense
        :value="value"
        @input="onInput"
        filled
        :items="typeItems"
        :rules="rules"
    >
    </v-select>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { SupportedParamType } from '../../../../ts/code'
import {
    TCodeParameterTypeMetadataOutput, getCodeParameterTypeMetadata
} from '../../../../ts/api/apiMetadata'
import { contactUsUrl } from '../../../../ts/url'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => null as SupportedParamType | null
        },
        rules: {
            type: Array,
            default: () => []
        },
        initialId: {
            type: Number,
            default: -1,
        },
    }
})

@Component
export default class SupportedParamTypeSelectComponent extends Props {
    supportedTypes : SupportedParamType[] = []

    get typeItems() : any[] {
        return this.supportedTypes.map((ele :SupportedParamType) => ({
            text: ele.Name,
            value: ele,
        }))
    }

    onInput(val : SupportedParamType) {
        this.$emit('input', val)
    }

    refreshSupportedTypes() {
        getCodeParameterTypeMetadata().then((resp : TCodeParameterTypeMetadataOutput) => {
            this.supportedTypes = resp.data
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
        this.refreshSupportedTypes()
    }

    @Watch('supportedTypes')
    @Watch('initialId')
    resyncFromId() {
        if (this.initialId == -1) {
            return
        }
        let idx : number = this.supportedTypes.findIndex((ele : SupportedParamType) => ele.Id == this.initialId)
        if (idx == -1) {
            return
        }
        this.onInput(this.supportedTypes[idx])
    }
}

</script>
