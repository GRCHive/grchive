<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Vendor
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-text-field v-model="url"
                      label="Url"
                      filled
                      :disabled="!canEdit">
        </v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :disabled="!canEdit">
        </v-textarea> 
    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit"
        >
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="success"
            @click="save"
            :disabled="!formValid"
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
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { Vendor } from '../../../ts/vendors'
import { updateVendor, newVendor, TNewVendorOutput } from '../../../ts/api/apiVendors'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        referenceVendor: {
            type: Object as () => Vendor | null,
            default: null
        }
    }
})

@Component
export default class CreateNewVendorForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name : string = ""
    url : string = ""
    description : string = ""

    doSave() {
        newVendor({
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            url : this.url,
            description: this.description,
        }).then((resp : TNewVendorOutput) => {
            this.$emit('do-save', resp.data)
            this.clearForm()
        }).catch((err: any) => {
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
        updateVendor({
            vendorId: this.referenceVendor!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            url : this.url,
            description: this.description,
        }).then((resp : TNewVendorOutput) => {
            this.$emit('do-save', resp.data)
            this.canEdit = false
        }).catch((err: any) => {
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

    cancel() {
        this.$emit('do-cancel')
        if (this.editMode) {
            this.canEdit = false
        }
    }

    edit() {
        this.canEdit = true
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }

    clearForm() {
        if (!!this.referenceVendor) {
            this.name = this.referenceVendor!.Name
            this.url = this.referenceVendor!.Url
            this.description = this.referenceVendor!.Description
        } else {
            this.name = ""
            this.url = ""
            this.description = ""
        }
    }
}

</script>
