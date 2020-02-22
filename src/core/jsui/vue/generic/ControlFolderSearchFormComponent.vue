<template>
    <v-autocomplete
        :value="value"
        @input="onInput"
        label="Folder"
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
import { contactUsUrl } from '../../ts/url'
import { FileFolder } from '../../ts/folders'
import { allControlFolderLink, TAllControlFolderLinkOutput } from '../../ts/api/apiControlFolderLinks'

const Props = Vue.extend({
    props: {
        controlId: Number,
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
export default class ControlFolderFormComponent extends Props {
    allFolders : FileFolder[] = []

    get items() : any[] {
        return this.allFolders.map((ele : FileFolder) => ({
            text: ele.Name,
            value: ele
        }))
    }

    refreshData() {
        allControlFolderLink({
            controlId: this.controlId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlFolderLinkOutput) => {
            this.allFolders = resp.data.Folders!
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
    
    onInput(c : FileFolder) {
        this.$emit('input', c)
    }

    mounted() {
        this.refreshData()
    }
}

</script>
