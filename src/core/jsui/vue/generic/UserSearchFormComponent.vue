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
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import { createUserString } from '../../ts/users'
import { getAllOrgUsers } from '../../ts/api/apiUsers'
import { contactUsUrl } from '../../ts/url'

export default Vue.extend({
    props: {
        label: String
    },
    data : () => ({
        loading: false,
        visibleUsers: [] as User[]
    }),
    computed: {
        displayItems() : string[] {
            let displayText = [] as string[]
            for (let user of this.visibleUsers) {
                displayText.push(createUserString(user))
            }
            return displayText
        }
    },
    methods: {
        loadAvailableUsers() {
            // TODO: Put this in metadata?
            this.loading = true
            getAllOrgUsers(<TGetAllOrgUsersInput>{
                //@ts-ignore
                csrf: this.$root.csrf,
                //@ts-ignore
                org: this.$root.orgGroupId
            }).then((resp : TGetAllOrgUsersOutput) => {
                this.visibleUsers = resp.data
                this.loading = false
            }).catch((err) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Please refresh the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    },
    mounted() {
        this.loadAvailableUsers()
    },
})
</script>
