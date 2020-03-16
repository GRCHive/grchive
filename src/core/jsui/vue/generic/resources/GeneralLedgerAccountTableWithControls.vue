<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    General Ledger Accounts
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
        </v-list-item>

        <v-divider></v-divider>

        <general-ledger-accounts-table
            :value="value"
            :resources="filteredAccounts"
            :search="filterText"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        >
        </general-ledger-accounts-table>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GeneralLedgerAccountsTable from '../GeneralLedgerAccountsTable.vue'
import { GeneralLedger, GeneralLedgerAccount } from '../../../ts/generalLedger'
import { TGetGLOutputs, getGL } from '../../../ts/api/apiGeneralLedger'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        value: {
            type: Array,
            default: () => [],
        },
        exclude: {
            type: Array,
            default: () => [],
        },
        enableSelect: {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        GeneralLedgerAccountsTable
    }
})
export default class GeneralLedgerAccountTableWithControls extends Props {
    filterText : string = ""
    showHideNew: boolean = false
    gl : GeneralLedger | null = null

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredAccounts() : GeneralLedgerAccount[] {
        if (!this.gl) {
            return []
        }
        return this.gl.listAccounts.filter((ele : GeneralLedgerAccount) => !this.excludeSet.has(ele.Id))
    }

    refreshGL() {
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

    modifySelected(vals : ProcessFlowControl[]) {
        this.$emit('input', vals)
    }

    mounted() {
        this.refreshGL()
    }
}

</script>
