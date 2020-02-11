<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} General Ledger {{ isSubledger ? "Subledger" : "Category" }}
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
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
            v-if="isSubledger || editMode"
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
            :items="finalAvailableCats"
            :rules="(!editMode && isSubledger) ? [rules.required] : []"
        ></v-autocomplete>
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
import * as rules from '../../../ts/formRules'
import { TNewGLCategoryInputs, TNewGLCategoryOutputs, newGLCategory } from '../../../ts/api/apiGeneralLedger'
import { TEditGLCategoryInputs, TEditGLCategoryOutputs, editGLCategory } from '../../../ts/api/apiGeneralLedger'
import { contactUsUrl } from '../../../ts/url'
import {PageParamsStore } from '../../../ts/pageParams'
import { GeneralLedgerCategory } from '../../../ts/generalLedger'

const VueComponent = Vue.extend({
    props: {
        isSubledger: Boolean,
        editMode: {
            type: Boolean,
            default: false
        },
        availableGlCats : {
            type: Array,
            default: () => []
        },
        referenceCat : {
            type: Object as () => GeneralLedgerCategory | null,
            default: null
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
    }
})

@Component
export default class CreateNewGeneralLedgerCategoryForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name: string = ""
    description: string = ""
    parentCategoryId: number | null = null

    get finalAvailableCats() : GeneralLedgerCategory[] {
        if (!this.referenceCat) {
            return this.availableGlCats as GeneralLedgerCategory[]
        }
        return (this.availableGlCats as GeneralLedgerCategory[]).filter((ele : GeneralLedgerCategory) =>
            ele.Id != this.referenceCat!.Id)
    }

    cancel() {
        this.$emit('do-cancel')

        if (this.editMode) {
            this.canEdit = false
        }
    }

    doSave() {
        newGLCategory(<TNewGLCategoryInputs>{
            orgId: PageParamsStore.state.organization!.Id,
            parentCategoryId: this.parentCategoryId,
            name: this.name,
            description: this.description
        }).then((resp : TNewGLCategoryOutputs) => {
            this.$emit('do-save', resp.data)
            this.resetForm()
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
        editGLCategory(<TEditGLCategoryInputs>{
            catId: this.referenceCat!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            parentCategoryId: this.parentCategoryId,
            name: this.name,
            description: this.description
        }).then((resp : TEditGLCategoryOutputs) => {
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

    resetForm() {
        if (!!this.referenceCat) {
            this.name = this.referenceCat.Name
            this.description = this.referenceCat.Description
            this.parentCategoryId = this.referenceCat.ParentCategoryId
        } else {
            this.name = ""
            this.description = ""
            this.parentCategoryId = null
        }
    }

    get canSubmit() : boolean {
        return this.formValid && this.name.length > 0
    }
}

</script>
