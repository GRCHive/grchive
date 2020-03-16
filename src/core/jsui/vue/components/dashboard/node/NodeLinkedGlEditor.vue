<template>
    <div>
        <v-list dense class="pa-0">
            <v-list-item class="pa-0">
                <v-list-item-action class="ma-0">
                    <v-btn icon @click="minimizeGL = !minimizeGL">
                        <v-icon small>
                            {{ !minimizeGL ? "mdi-chevron-up" : "mdi-chevron-down" }}
                        </v-icon>
                    </v-btn>
                </v-list-item-action>

                <v-subheader class="flex-grow-1 pr-0">
                    LINKED GENERAL LEDGER ACCOUNTS
                </v-subheader>

                <v-list-item-action class="ma-0">
                    <v-dialog persistent max-width="40%" v-model="showLinkGL">
                        <template v-slot:activator="{ on }">
                            <v-btn
                                icon
                                v-on="on"
                            >
                                <v-icon small>
                                    mdi-plus
                                </v-icon>
                            </v-btn>
                        </template>

                        <v-card>
                            <v-card-title>
                                Link General Ledger Accounts
                            </v-card-title>

                            <general-ledger-account-search-form-component
                                v-model="accountsToLink"
                            >
                            </general-ledger-account-search-form-component>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="cancelGLLink"
                                >
                                    Cancel
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-btn
                                    color="success"
                                    @click="saveGLLink"
                                    :disabled="accountsToLink.length == 0"
                                >
                                    Link
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </v-list-item-action>
            </v-list-item>
        </v-list>
        <general-ledger-accounts-table
            :resources="linkedGL"
            use-crud-delete
            @delete="deleteLinkedGL"
            v-if="!minimizeGL"
        >
        </general-ledger-accounts-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import VueSetup from '../../../../ts/vueSetup' 
import GeneralLedgerAccountSearchFormComponent from '../../../generic/GeneralLedgerAccountSearchFormComponent.vue'
import GeneralLedgerAccountsTable from '../../../generic/GeneralLedgerAccountsTable.vue'
import { PageParamsStore } from '../../../../ts/pageParams'
import { contactUsUrl } from '../../../../ts/url'
import { GeneralLedgerAccount, GeneralLedger } from '../../../../ts/generalLedger'
import { 
    newNodeGLLink,
    deleteNodeGLLink
} from '../../../../ts/api/apiNodeGLLinks'

@Component({
    components: {
        GeneralLedgerAccountsTable,
        GeneralLedgerAccountSearchFormComponent
    }
})
export default class NodeLinkedGLEditor extends Vue {
    accountsToLink: GeneralLedgerAccount[] = []
    showLinkGL: boolean = false
    minimizeGL: boolean = false

    cancelGLLink() {
        this.accountsToLink = []
        this.showLinkGL = false
    }

    saveGLLink() {
        if (this.accountsToLink.length == 0) {
            return
        }

        let nodeId : number = this.currentNode.Id
        let account : GeneralLedgerAccount = this.accountsToLink[0]
        newNodeGLLink({
            nodeId: nodeId,
            accountId: account.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            VueSetup.store.commit('addNodeGLLink', {
                nodeId: nodeId,
                account: account,
            })
            this.accountsToLink = []
            this.showLinkGL = false
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

    deleteLinkedGL(account : GeneralLedgerAccount) {
        let id : number = this.currentNode.Id
        deleteNodeGLLink({
            nodeId: id,
            accountId: account.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            VueSetup.store.commit('deleteNodeGLLink', {
                nodeId: id,
                accountId: account.Id,
            })
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

    get linkedGL(): GeneralLedgerAccount[] | null {
        let gl : GeneralLedger | null = VueSetup.store.getters.glLinkedToNode(this.currentNode.Id)
        if (!gl) {
            return null
        }
        return gl.changed && gl.listAccounts
    }

    get currentNode() : ProcessFlowNode {
        return VueSetup.store.getters.currentNodeInfo
    }
}

</script>
