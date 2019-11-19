<template>
    <div>
        <v-dialog
            v-if="!!currentEditDeleteCat"
            v-model="showHideEditCat"
            persistent
            max-width="40%">

            <create-new-general-ledger-category-form
                :edit-mode="true"
                :is-subledger="!!currentEditDeleteCat.ParentCategory"
                :available-gl-cats="availableCats"
                :reference-cat="currentEditDeleteCat"
                ref="editForm"
                @do-cancel="cancelCatEdit"
                @do-save="finishCatEdit"
            ></create-new-general-ledger-category-form>
        </v-dialog>

        <v-dialog
            v-if="!!currentEditDeleteCat"
            v-model="showHideDeleteCat"
            persistent
            max-width="40%">

            <generic-delete-confirmation-form
                :item-name="!!currentEditDeleteCat.ParentCategory ? 'subledgers' : 'categories'"
                :itemsToDelete="[currentEditDeleteCat.Name]"
                :use-global-deletion="false"
                @do-cancel="cancelCatDelete"
                @do-delete="confirmCatDelete"
            >
            </generic-delete-confirmation-form>
        </v-dialog>

        <v-treeview
            ref="view"
            :items="displayItems"
            open-on-click
            open-all>
            <template v-slot:prepend="{ item, open }">
                <v-icon v-if="!!item.cat">
                    {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                </v-icon>
            </template>

            <template v-slot:label="{ item }">
                <span :class="!!item.cat ? 'font-weight-bold' : ''">{{ item.name }}</span>
            </template>

            <template v-slot:append="{ item }">
                <template v-if="!!item.cat">
                    <v-btn icon @click.stop="onEditCat(item.cat)">
                        <v-icon>mdi-pencil</v-icon>
                    </v-btn>

                    <v-btn icon @click.stop="onDeleteCat(item.cat)">
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </template>
            </template>
        </v-treeview>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { RawGeneralLedgerCategory,
         RawGeneralLedgerAccount,
         GeneralLedgerCategory,
         GeneralLedgerAccount,
         GeneralLedger } from '../../ts/generalLedger'
import { getGL, TGetGLInputs, TGetGLOutputs } from '../../ts/api/apiGeneralLedger'
import { deleteGLCategory, TDeleteGLCategoryInputs, TDeleteGLCategoryOutputs } from '../../ts/api/apiGeneralLedger'
import { contactUsUrl } from '../../ts/url'
import CreateNewGeneralLedgerCategoryForm from '../components/dashboard/CreateNewGeneralLedgerCategoryForm.vue'
import GenericDeleteConfirmationForm from '../components/dashboard/GenericDeleteConfirmationForm.vue'

const VueComponent = Vue.extend({
    props: {
        orgId: Number 
    },
    components: {
        CreateNewGeneralLedgerCategoryForm,
        GenericDeleteConfirmationForm
    },
})

@Component
export default class GeneralLedgerDisplay extends VueComponent {
    generalLedger: GeneralLedger = new GeneralLedger()
    showHideEditCat: boolean = false
    showHideDeleteCat: boolean = false
    currentEditDeleteCat : GeneralLedgerCategory | null = null

    $refs!: {
        view: any,
        editForm: CreateNewGeneralLedgerCategoryForm
    }

    get availableCats() : GeneralLedgerCategory[] {
        return this.generalLedger.changed && this.generalLedger.listCategories
    }

    refreshGeneralLedger() {
        getGL({
            orgId: this.orgId
        }).then((resp : TGetGLOutputs) => {
            this.generalLedger.rebuildGL(resp.data.Categories, resp.data.Accounts)
            Vue.nextTick(() => {
                this.$refs.view.updateAll(true)
            })
        }).catch((err : any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    addCategory(cat : RawGeneralLedgerCategory) {
        this.generalLedger.addRawCategory(cat)
    }

    addAccount(acc : RawGeneralLedgerAccount) {
        this.generalLedger.addRawAccount(acc)
    }

    mounted() {
        this.refreshGeneralLedger()
    }

    createDisplayItemForCategory(cat : GeneralLedgerCategory) : any {
        let childCats = [] as any[]
        let childAccs = [] as any[]

        for (let id of cat.SubCategories.keys()) {
            childCats.push(this.createDisplayItemForCategory(cat.SubCategories.get(id)!))
        }

        for (let id of cat.SubAccounts.keys()) {
            childAccs.push(this.createDisplayItemForAccount(cat.SubAccounts.get(id)!))
        }

        return {
            id: cat.Id.toString(),
            name: cat.Name,
            cat: cat,
            children: [...childCats, ...childAccs]
        }
    }

    createDisplayItemForAccount(acc : RawGeneralLedgerAccount) : any {
        return {
            id: `${acc.ParentCategoryId} ${acc.Id}`,
            name: `${acc.AccountName} ${acc.AccountName != acc.AccountId ?  "(" + acc.AccountId.toString() + ")" : ""}`,
            acc: acc
        }
    }

    get displayItems() : any[] {
        let items = [] as any[]
        for (let id of this.generalLedger.topLevelCategories.keys()) {
            items.push(this.createDisplayItemForCategory(this.generalLedger.topLevelCategories.get(id)!))
        }
        return this.generalLedger.changed && items
    }

    onEditCat(cat : GeneralLedgerCategory) {
        this.currentEditDeleteCat = cat
        this.showHideEditCat = true
        Vue.nextTick(() => {
            this.$refs.editForm.resetForm()
        })
    }

    onDeleteCat(cat : GeneralLedgerCategory) {
        this.currentEditDeleteCat = cat
        this.showHideDeleteCat = true
    }

    cancelCatDelete() {
        this.currentEditDeleteCat = null
        this.showHideDeleteCat = false
    }

    confirmCatDelete() {
        deleteGLCategory(<TDeleteGLCategoryInputs>{
            orgId: this.orgId,
            catId: this.currentEditDeleteCat!.Id,
        }).then((resp : TDeleteGLCategoryOutputs) => {
            this.generalLedger.removeCategory(this.currentEditDeleteCat!.Id)
            this.currentEditDeleteCat = null
            this.showHideDeleteCat = false
        }).catch((err: any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    cancelCatEdit() {
        this.showHideEditCat = false
    }

    finishCatEdit(cat : RawGeneralLedgerCategory) {
        this.showHideEditCat = false
        this.generalLedger.replaceRawCategory(cat)
    }
}

</script>
