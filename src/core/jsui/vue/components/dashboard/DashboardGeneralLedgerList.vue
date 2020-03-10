<template>
    <div class="ma-4">
        <v-dialog
            v-model="showHideNewCategory"
            persistent
            max-width="40%"
        >
            <create-new-general-ledger-category-form
                :is-subledger="false"
                @do-save="onCreateNewCategory"
                @do-cancel="onCancelNewCategory">
            </create-new-general-ledger-category-form>
        </v-dialog>

        <v-dialog
            v-model="showHideNewSubledger"
            persistent
            max-width="40%"
        >
            <create-new-general-ledger-category-form
                :is-subledger="true"
                :available-gl-cats="availableGLCats"
                @do-save="onCreateNewSubledger"
                @do-cancel="onCancelNewSubledger">
            </create-new-general-ledger-category-form>
        </v-dialog>

        <v-dialog
            v-model="showHideNewAccount"
            persistent
            max-width="40%"
        >
            <create-new-general-ledger-account-form
                :available-gl-cats="availableGLCats"
                @do-save="onCreateNewAccount"
                @do-cancel="onCancelNewAccount">
            </create-new-general-ledger-account-form>
        </v-dialog>

        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    General Ledger
                </v-list-item-title>
            </v-list-item-content>

            <v-spacer></v-spacer>
            <v-list-item-action>
                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <v-list dense>
                        <v-list-item @click.stop="showHideNewCategory = true">
                            <v-list-item-title>
                                Category
                            </v-list-item-title>
                        </v-list-item>
                        <v-list-item @click.stop="showHideNewSubledger = true">
                            <v-list-item-title>
                                Subledger
                            </v-list-item-title>
                        </v-list-item>
                        <v-list-item @click.stop="showHideNewAccount = true">
                            <v-list-item-title>
                                Account
                            </v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <v-tabs>
            <v-tab>Ledger</v-tab>
            <v-tab-item>
                <general-ledger-display
                    :org-id="orgId"
                    ref="ledger"
                ></general-ledger-display>
            </v-tab-item>

            <v-tab>Audit Trail</v-tab>
            <v-tab-item>
            </v-tab-item>
        </v-tabs>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GeneralLedgerDisplay from '../../generic/GeneralLedgerDisplay.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import { GeneralLedgerCategory, RawGeneralLedgerCategory, RawGeneralLedgerAccount } from '../../../ts/generalLedger'
import CreateNewGeneralLedgerCategoryForm from './CreateNewGeneralLedgerCategoryForm.vue'
import CreateNewGeneralLedgerAccountForm from './CreateNewGeneralLedgerAccountForm.vue'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'

@Component({
    components: {
        GeneralLedgerDisplay,
        CreateNewGeneralLedgerAccountForm,
        CreateNewGeneralLedgerCategoryForm,
        AuditTrailViewer
    }
})
export default class DashboardGeneralLedgerList extends Vue {
    showHideNewCategory : boolean = false
    showHideNewSubledger : boolean = false
    showHideNewAccount : boolean = false
    isMounted: boolean = false

    $refs!: {
        ledger: GeneralLedgerDisplay
    }

    get orgId() : number {
        return PageParamsStore.state.organization!.Id
    }

    get availableGLCats() : GeneralLedgerCategory[] {
        if (!this.isMounted) {
            return []
        }

        return this.$refs.ledger.generalLedger.changed && 
            this.$refs.ledger.generalLedger.listCategories
    }

    onCancelNewCategory() {
        this.showHideNewCategory = false
    }

    onCreateNewCategory(cat : RawGeneralLedgerCategory) {
        this.$refs.ledger.addCategory(cat)
        this.showHideNewCategory = false
    }

    onCancelNewSubledger() {
        this.showHideNewSubledger = false
    }

    onCreateNewSubledger(cat : RawGeneralLedgerCategory) {
        this.$refs.ledger.addCategory(cat)
        this.showHideNewSubledger = false
    }

    onCancelNewAccount() {
        this.showHideNewAccount = false
    }

    onCreateNewAccount(acc : RawGeneralLedgerAccount) {
        this.$refs.ledger.addAccount(acc)
        this.showHideNewAccount = false
    }

    mounted() {
        this.isMounted = true
    }
}

</script>
