<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content>
                <v-list-item-title class="title">
                    My Organizations
                </v-list-item-title>
            </v-list-item-content>
        </v-list-item>
        <v-divider></v-divider>

        <v-card
            v-for="(item, index) in availableOrgs"
            :key="index"
            class="my-2"
        >
            <v-list-item :href="createOrgUrl(item.OktaGroupName)">
                <v-list-item-content>
                    <v-list-item-title>
                        {{ item.Name }}
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>
        </v-card>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { TGetUserOrgsInput, TGetUserOrgsOutput, getAllOrgsForUser} from '../../../ts/api/apiUsers'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgUrl } from '../../../ts/url'
import { Organization } from '../../../ts/organizations'

export default Vue.extend({
    data: () => ({
        availableOrgs: [] as Organization[]
    }),
    mounted() {
        getAllOrgsForUser(<TGetUserOrgsInput>{
            userId: PageParamsStore.state.user!.Id
        }).then((resp : TGetUserOrgsOutput) => {
            this.availableOrgs = resp.data
        }).catch((err: any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    },
    methods: {
        createOrgUrl
    }
})

</script>
