<template>
    <v-autocomplete
        label="Vendor"
        filled
        cache-items
        :items="displayItems"
        :loading="loading"
        hide-no-data
        hide-selected
        clearable
        :disabled="disabled"
        :rules="rules"
        :value="value"
        @input="onChangeVendor"
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Vendor } from '../../ts/vendors'
import { PageParamsStore } from '../../ts/pageParams'
import { allVendors, TAllVendorOutput } from '../../ts/api/apiVendors'
import { contactUsUrl } from '../../ts/url'

const Props = Vue.extend({
    props: {
        rules: {
            type: Array,
            default: []
        },
        disabled: {
            type: Boolean,
            default: false
        },
        value: {
            type: Object as () => Vendor | null,
            default: null
        },
        initialVendorId: {
            type: Number,
            default: -1
        }
    }
})

@Component
export default class VendorSearchFormComponent extends Props {
    loading : boolean = false
    allVendors : Vendor[] = []

    get displayItems() : any[] {
        return this.allVendors.map((ele : Vendor) => ({
            text: ele.Name,
            value: ele
        }))
    }

    onChangeVendor(v : Vendor) {
        this.$emit('input', v)
    }

    refreshData() {
        this.loading = true
        allVendors({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllVendorOutput) => {
            this.allVendors = resp.data
            this.loading = false

            if (this.initialVendorId >= 0) {
                let idx = this.allVendors.findIndex((ele : Vendor) => ele.Id == this.initialVendorId)
                if (idx != -1) {
                    this.onChangeVendor(this.allVendors[idx])
                }
            }
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

    mounted() {
        this.refreshData()
    }
}

</script>
