<template>
    <general-ledger-accounts-table
        :resources="accounts"
        :value="value"
        @input="onInput"
        selectable
    >
    </general-ledger-accounts-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { GeneralLedger, GeneralLedgerAccount } from '../../ts/generalLedger'
import { TGetGLOutputs, getGL } from '../../ts/api/apiGeneralLedger'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import GeneralLedgerAccountsTable from './GeneralLedgerAccountsTable.vue'

const Props = Vue.extend({
    props: {
        value: Array,
    }
})

@Component({
    components: {
        GeneralLedgerAccountsTable
    }
})
export default class GeneralLedgerAccountSearchFormComponent extends Props {
    gl : GeneralLedger | null = null

    get accounts() : GeneralLedgerAccount[] {
        if (!this.gl) {
            return []
        }
        return this.gl.listAccounts
    }

    onInput(v : GeneralLedgerAccount[]) {
        this.$emit('input', v)
    }

    refreshData() {
        getGL({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetGLOutputs) => {
            this.gl = Vue.observable(new GeneralLedger())
            this.gl.rebuildGL(resp.data.Categories, resp.data.Accounts)
        }).catch((err : any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>
