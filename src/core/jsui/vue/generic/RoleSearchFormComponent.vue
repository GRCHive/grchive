<template>
    <v-autocomplete
        :label="label"
        filled
        cache-items
        :items="displayItems"
        :loading="loading"
        hide-no-data
        hide-selected
        clearable
        :value="role"
        @change="changeRole"
        :disabled="disabled"
        :rules="rules"
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { RoleMetadata } from '../../ts/roles'
import { TGetAllOrgRolesOutput, getAllOrgRoles } from '../../ts/api/apiRoles'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'

const Props = Vue.extend({
    props: {
        label: String,
        role: {
            type: Object as () => RoleMetadata | null
        },
        disabled: {
            type: Boolean,
            default: false
        },
        preloadRoles: {
            type: Object as () => RoleMetadata[] | null,
            default: null
        },
        rules: Array,
    },
})

@Component
export default class RoleSearchFormComponent extends Props {
    loading: boolean = false
    availableRoles: RoleMetadata[] = []

    get displayItems() : Object[] {
        return this.availableRoles.map((ele : RoleMetadata) => ({
            text: ele.Name,
            value: ele
        }))
    }

    changeRole(val : RoleMetadata) {
        this.$emit('update:role', val)
    }

    mounted() {
        if (!!this.preloadRoles) {
            this.availableRoles = this.preloadRoles!
            this.loading = false
        } else {
            getAllOrgRoles({
                orgId: PageParamsStore.state.organization!.Id
            }).then((resp : TGetAllOrgRolesOutput) => {
                this.availableRoles = resp.data
                this.loading = false
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
    }
}

</script>
