<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} General Ledger Account
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="id"
                      label="Identifier"
                      filled
                      :rules="[rules.required]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :disabled="!canEdit">
        </v-textarea> 

        <v-autocomplete
            v-model="parentCategoryId"
            label="Parent Category"
            deletable-chips
            chips
            clearable
            :disabled="!canEdit"
            hide-no-data
            hide-selected
            filled
            item-text="Name"
            item-value="Id"
            :items="availableGlCats"
            :rules="[rules.required]"
        ></v-autocomplete>

        <v-checkbox
            label="Financially Relevant"
            v-model="financiallyRelevant">
        </v-checkbox>
    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit || dialogMode"
        >
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="success"
            @click="save"
            :disabled="!canSubmit"
            v-if="canEdit"
        >
            Save
        </v-btn>

        <v-btn
            color="success"
            @click="edit"
            v-if="!canEdit"
        >
            Edit
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import * as rules from '../../../ts/formRules'
import { GeneralLedgerAccount } from '../../../ts/generalLedger'
import { TNewGLAccountInputs, TNewGLAccountOutputs, newGLAccount } from '../../../ts/api/apiGeneralLedger'
import { TEditGLAccountInputs, TEditGLAccountOutputs, editGLAccount } from '../../../ts/api/apiGeneralLedger'
import { contactUsUrl } from '../../../ts/url'
import {PageParamsStore } from '../../../ts/pageParams'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        availableGlCats : {
            type: Array,
            default: () => []
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
        referenceAccount: {
            type: Object as () => GeneralLedgerAccount | null,
            default: null
        }
    }
})

@Component
export default class CreateNewGeneralLedgerAccountForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    parentCategoryId: number | null = null
    id: string = ""
    name: string = ""
    description: string = ""
    financiallyRelevant: boolean = true

    cancel() {
        this.$emit('do-cancel')

        if (this.editMode) {
            this.canEdit = false
        }
    }

    doSave() {
        newGLAccount(<TNewGLAccountInputs>{
            orgId: PageParamsStore.state.organization!.Id,
            parentCategoryId: this.parentCategoryId!,
            accountId: this.id,
            accountName: this.name,
            accountDescription: this.description,
            financiallyRelevant: this.financiallyRelevant
        }).then((resp : TNewGLAccountOutputs) => {
            this.$emit('do-save', resp.data)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    doEdit() {
        editGLAccount(<TEditGLAccountInputs>{
            accId: this.referenceAccount!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            parentCategoryId: this.parentCategoryId!,
            accountId: this.id,
            accountName: this.name,
            accountDescription: this.description,
            financiallyRelevant: this.financiallyRelevant
        }).then((resp : TEditGLAccountOutputs) => {
            this.canEdit = false
            this.$emit('do-save', resp.data)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    save() {
        if (this.editMode) {
            this.doEdit()
        } else {
            this.doSave()
        }
    }

    edit() {
        this.canEdit = true
    }

    mounted() {
        this.canEdit = !this.editMode
    }

    get canSubmit() : boolean {
        return this.formValid
    }

    @Watch('id')
    onIdChange(newId : string, oldId : string) {
        if (oldId != this.name && this.name.length > 0) {
            return
        }
        this.name = newId
    }

    resetForm() {
        if (!!this.referenceAccount) {
            this.parentCategoryId = this.referenceAccount.ParentCategoryId
            this.id = this.referenceAccount.AccountId
            this.name = this.referenceAccount.AccountName
            this.description = this.referenceAccount.AccountDescription
            this.financiallyRelevant = this.referenceAccount.FinanciallyRelevant
        } else {
            this.parentCategoryId = null
            this.id = ""
            this.name = ""
            this.description = ""
            this.financiallyRelevant = true
        }
    }
}


</script>
