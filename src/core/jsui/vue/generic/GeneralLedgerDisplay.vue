<template>
    <v-treeview
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
    </v-treeview>
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
import { contactUsUrl } from '../../ts/url'

const VueComponent = Vue.extend({
    props: {
        orgId: Number 
    }
})

@Component
export default class GeneralLedgerDisplay extends VueComponent {
    generalLedger: GeneralLedger = new GeneralLedger()

    refreshGeneralLedger() {
        getGL({
            orgId: this.orgId
        }).then((resp : TGetGLOutputs) => {
            this.generalLedger.rebuildGL(resp.data.Categories, resp.data.Accounts)
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
}

</script>
