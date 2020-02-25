<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-breadcrumbs :items="parentBreadcrumbs" class="pa-0">
            </v-breadcrumbs>
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Account: {{ glAccount.AccountName }}
                        <span class="subtitle-1" v-if="glAccount.AccountName != glAccount.AccountId">
                            ({{ glAccount.AccountId }})
                        </span>
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ glAccount.AccountDescription }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-spacer></v-spacer>

                <v-dialog v-model="showHideDelete"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="error" v-on="on">
                            Delete
                        </v-btn>
                    </template>

                    <generic-delete-confirmation-form
                        item-name="accounts"
                        :items-to-delete="[glAccount.AccountName]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="5">
                        <create-new-general-ledger-account-form
                            :edit-mode="true"
                            :reference-account="glAccount"
                            :available-gl-cats="availableGLCats"
                            @do-save="finishEdit"
                            ref="editForm">
                        </create-new-general-ledger-account-form>
                    </v-col>

                    <v-col cols="7">
                        <v-card>
                            <v-card-title>
                                Related Resources
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-tabs>
                                <v-tab>Process Flows</v-tab>
                                <v-tab-item>
                                    <process-flow-table
                                        :resources="relatedFlows"
                                    >
                                    </process-flow-table>
                                </v-tab-item>

                                <v-tab>Risks</v-tab>
                                <v-tab-item>
                                    <risk-table
                                        :resources="relatedRisks"
                                    >
                                    </risk-table>
                                </v-tab-item>

                                <v-tab>Controls</v-tab>
                                <v-tab-item>
                                    <control-table
                                        :resources="relatedControls"
                                    >
                                    </control-table>
                                </v-tab-item>
                            </v-tabs>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { RawGeneralLedgerAccount, GeneralLedgerAccount, GeneralLedgerCategory, GeneralLedger } from '../../../ts/generalLedger'
import { TGetGLAccountInputs, TGetGLAccountOutputs, getGLAccount } from '../../../ts/api/apiGeneralLedger'
import { TDeleteGLAccountInputs, TDeleteGLAccountOutputs, deleteGLAccount } from '../../../ts/api/apiGeneralLedger'
import { PageParamsStore } from '../../../ts/pageParams'
import { createOrgGLUrl, contactUsUrl } from '../../../ts/url'
import CreateNewGeneralLedgerAccountForm from './CreateNewGeneralLedgerAccountForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import RiskTable from '../../generic/RiskTable.vue'
import ControlTable from '../../generic/ControlTable.vue'
import ProcessFlowTable from '../../generic/ProcessFlowTable.vue'
import { TAllRiskGLLinkOutput, allRiskGLLink } from '../../../ts/api/apiRiskGLLinks'
import { TAllControlGLLinkOutput, allControlGLLink } from '../../../ts/api/apiControlGLLinks'
import { TAllNodeGLLinkOutput, allNodeGLLink } from '../../../ts/api/apiNodeGLLinks'

@Component({
    components: {
        CreateNewGeneralLedgerAccountForm,
        GenericDeleteConfirmationForm,
        RiskTable,
        ControlTable,
        ProcessFlowTable
    }
})
export default class FullEditGeneralLedgerAccountComponent extends Vue {
    ready: boolean = false
    ledger : GeneralLedger = new GeneralLedger()
    expandDescription : boolean = false
    showHideDelete: boolean = false

    relatedControls : ProcessFlowControl[] = []
    relatedRisks : ProcessFlowRisk[] = []
    relatedFlows : ProcessFlowBasicData[] = []

    $refs!: {
        editForm: CreateNewGeneralLedgerAccountForm
    }

    get parentBreadcrumbs() : any[] {
        let parentCrumbs = []
        let currentParent : GeneralLedgerCategory | null = this.glAccount.ParentCategory
        while (currentParent != null) {
            parentCrumbs.unshift({
                disabled: true,
                text: currentParent.Name
            })
            currentParent = currentParent.ParentCategory
        }
        return this.ledger.changed && parentCrumbs
    }

    get glAccount() : GeneralLedgerAccount {
        return this.ledger.accounts.values().next().value
    }

    get availableGLCats() : GeneralLedgerCategory[] {
        return this.ledger.changed && 
            this.ledger.listCategories
    }

    refreshRelatedFlows() {
        allNodeGLLink({
            accountId: this.glAccount.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllNodeGLLinkOutput) => {
            this.relatedFlows = resp.data.Flows!
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

    refreshRelatedRisks() {
        allRiskGLLink({
            accountId: this.glAccount.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllRiskGLLinkOutput) => {
            this.relatedRisks = resp.data.Risks!
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

    refreshRelatedControls() {
        allControlGLLink({
            accountId: this.glAccount.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlGLLinkOutput) => {
            this.relatedControls = resp.data.Controls!
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

    refreshAccountData() {
        let data = window.location.pathname.split('/')
        let accId = Number(data[data.length - 1])

        getGLAccount(<TGetGLAccountInputs>{
            orgId: PageParamsStore.state.organization!.Id,
            accId: accId,
        }).then((resp : TGetGLAccountOutputs) => {
            this.ledger.rebuildGL(resp.data.Parents, [resp.data.Account])
            this.ready = true

            this.refreshRelatedFlows()
            this.refreshRelatedRisks()
            this.refreshRelatedControls()

            Vue.nextTick(() => {
                this.$refs.editForm.resetForm()
            })
        }).catch((err : any) => {
            window.location.replace('/404')
        })
    }

    mounted() {
        this.refreshAccountData()
    }

    onDelete() {
        deleteGLAccount(<TDeleteGLAccountInputs>{
            orgId: PageParamsStore.state.organization!.Id,
            accId: this.glAccount.Id,
        }).then((resp : TDeleteGLAccountOutputs) => {
            window.location.replace(createOrgGLUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    finishEdit(acc : RawGeneralLedgerAccount) {
        this.ledger.replaceRawAccount(acc)
    }
}

</script>
