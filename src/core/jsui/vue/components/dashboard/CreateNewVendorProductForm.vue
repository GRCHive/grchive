<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Vendor Product
    </v-card-title>

    <v-card-subtitle>
        Vendor: {{ parentVendor.Name }}
    </v-card-subtitle>
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
import { VendorProduct, Vendor } from '../../../ts/vendors'
import * as rules from '../../../ts/formRules'
import { contactUsUrl } from '../../../ts/url'
import { newVendorProduct, updateVendorProduct, TNewVendorProductOutput } from '../../../ts/api/apiVendorProduct'
import { PageParamsStore } from '../../../ts/pageParams'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        referenceProduct: {
            type: Object as () => VendorProduct | null,
            default: null
        },
        parentVendor: {
            type: Object as () => Vendor | null,
            required: true
        }
    }
})

@Component
export default class CreateNewVendorProductForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name : string = ""
    url : string = ""
    description : string = ""

    doSave() {
        newVendorProduct({
            orgId: PageParamsStore.state.organization!.Id,
            vendorId: this.parentVendor!.Id,
            name: this.name,
            url : this.url,
            description: this.description,
        }).then((resp : TNewVendorProductOutput) => {
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
        updateVendorProduct({
            productId: this.referenceProduct!.Id,
            vendorId: this.parentVendor!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            url : this.url,
            description: this.description,
        }).then((resp : TNewVendorProductOutput) => {
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
        if (!!this.referenceProduct) {
            this.name = this.referenceProduct!.Name
            this.url = this.referenceProduct!.Url
            this.description = this.referenceProduct!.Description
        } else {
            this.name = ""
            this.url = ""
            this.description = ""
        }
    }
}

</script>
