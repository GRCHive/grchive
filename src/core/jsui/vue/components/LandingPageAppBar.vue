<template>
    <v-app-bar>
        <a :href="homePageUrl">
            <v-img
                src="/static/assets/logos/grchive.png"
                height="40px"
                max-height="40px"
                max-width="254px"
                contain
            >
            </v-img>
        </a>
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
                :href="blogUrl"
            >
                Blog
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

import { contactUsUrl, loginPageUrl, homePageUrl, dashboardUrl, blogUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

export default {
    data: () => ({
        contactUsUrl,
        loginPageUrl,
        homePageUrl,
        dashboardUrl,
        blogUrl
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
