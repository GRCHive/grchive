<template>
    <v-app-bar app>
        <v-toolbar-title color="primary">
            <a :href="homePageUrl">{{ companyName }} </a>
        </v-toolbar-title>
        <div class="flex-grow-1"></div>

        <v-toolbar-items>
            <v-btn
                text color="primary"
                :href="contactUsUrl"
            >
                Contact Us
            </v-btn>

            <v-btn
                text color="primary"
                :href="dashboardUrl"
                v-if="hasAuth"
            >
                Dashboard
            </v-btn>

            <v-btn
                text color="primary"
                :href="loginPageUrl"
                v-else
            >
                Login
            </v-btn>

        </v-toolbar-items>
    </v-app-bar>
</template>

<script lang="ts">

import { contactUsUrl, loginPageUrl, homePageUrl, dashboardUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

export default {
    data: () => ({
        contactUsUrl,
        loginPageUrl,
        homePageUrl,
        dashboardUrl
    }),
    props: {
        companyName : String
    },
    computed: {
        hasAuth() : boolean {
            return PageParamsStore.state.user!.Auth
        }
    }
}

</script>

<style scoped>

a {
    text-decoration: none;
    color: black !important;
}

</style>
