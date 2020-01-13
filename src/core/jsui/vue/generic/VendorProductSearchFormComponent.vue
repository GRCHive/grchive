<template>
    <div>
        <vendor-search-form-component
            :rules="[importRules.required]"
            :disabled="disabled"
            v-model="parentVendor"
            :initial-vendor-id="initialVendorId"
        >
        </vendor-search-form-component>

        <v-autocomplete
            label="Product"
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
            @input="onChangeProduct"
        ></v-autocomplete>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import VendorSearchFormComponent from './VendorSearchFormComponent.vue'
import * as rules from '../../ts/formRules'
import { Vendor, VendorProduct } from '../../ts/vendors'
import { Watch } from 'vue-property-decorator'
import { allVendorProducts, TAllVendorProductOutput } from '../../ts/api/apiVendorProduct'
import { PageParamsStore } from '../../ts/pageParams'
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
            type: Object as () => VendorProduct | null,
            default: null
        },
        initialVendor: {
            type: Object as () => Vendor | null,
            default: null
        },
        initialVendorId: {
            type: Number,
            default: -1
        }
    }
})

@Component({
    components: {
        VendorSearchFormComponent,
    }
})
export default class VendorProductSearchFormComponent extends Props {
    importRules : any = rules
    parentVendor : Vendor | null = null

    loading : boolean = false
    availableProducts: VendorProduct[] | null = null

    get displayItems() : any[] {
        if (!this.availableProducts) {
            return []
        }
        return this.availableProducts.map((ele : VendorProduct) => ({
            text : ele.Name,
            value: ele
        }))
    }

    onChangeProduct(p : VendorProduct) {
        this.$emit('input', p)
    }

    @Watch('parentVendor')
    reloadProducts() {
        if (!this.parentVendor) {
            this.availableProducts = null
            return
        }

        allVendorProducts({
            orgId: PageParamsStore.state.organization!.Id,
            vendorId: this.parentVendor!.Id,
        }).then((resp : TAllVendorProductOutput) => {
            this.availableProducts = resp.data
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
        if (!!this.initialVendor) {
            this.parentVendor = this.initialVendor
        }
    }
}

</script>
