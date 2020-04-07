<template>
    <generic-code-ide
        :value="codeString"
        :lang="lang"
        :readonly="readonly"
        :full-height="fullHeight"
        :save-in-progress="saveInProgress"
        @input="onInput"
        @save="onSave"
    >
    </generic-code-ide>
</template>

<script lang="ts">

import Vue from 'vue'
import Component, { mixins } from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { Props } from './GenericCodeEditor.vue'
import GenericCodeIde from './GenericCodeIDE.vue'
import { ManagedCode } from '../../../ts/code'
import { PageParamsStore } from '../../../ts/pageParams'
import { 
    getCode, TGetCodeOutput,
    allCode, TAllCodeOutput,
    saveCode, TSaveCodeOutput,
} from '../../../ts/api/apiCode'
import { contactUsUrl } from '../../../ts/url'

const ManagedProps = Vue.extend({
    props: {
        dataId: {
            type: Number,
            default: -1,
        }
    }
})

@Component({
    components: {
        GenericCodeIde,
    }
})
export default class ManagedCodeIDE extends mixins(Props, ManagedProps) {
    codeString : string = ""
    loading: boolean = false
    saveInProgress: boolean = false

    allCode : ManagedCode[] = []
    selectedCode : ManagedCode | null = null

    onInput(text : string) {
        this.codeString = text
    }

    @Watch('selectedCode')
    pullCode() {
        if (!this.selectedCode) {
            this.codeString = ""
            return
        }

        this.loading = true
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
            codeId: this.selectedCode.Id,
        }

        if (this.dataId != -1) {
            params.dataId = this.dataId
        }

        getCode(params).then((resp : TGetCodeOutput) => {
            this.codeString = resp.data
            this.loading = false
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

    @Watch('dataId')
    refreshCode() {
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
        }

        if (this.dataId != -1) {
            params.dataId = this.dataId
        }

        allCode(params).then((resp : TAllCodeOutput) => {
            this.allCode = resp.data
            if (this.allCode.length > 0) {
                this.selectedCode = this.allCode[0]
            } else {
                this.selectedCode = null
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

    onSave() {
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
            code: this.codeString,
        }

        this.saveInProgress = true

        if (this.dataId != -1) {
            params.dataId = this.dataId
        }

        saveCode(params).then((resp : TSaveCodeOutput) => {
            this.allCode.unshift(resp.data)
            this.selectedCode = resp.data
            this.saveInProgress = false
        }).catch((err : any) => {
            this.saveInProgress = false
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
        this.refreshCode()
    }
}

</script>

