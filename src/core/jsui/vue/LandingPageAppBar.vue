<template>
    <v-app-bar app>
        <v-toolbar-title color="primary">
            <a :href="homePageUrl">{{ companyName }}</a>
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
                :href="loginPageUrl"
            >
                Login
            </v-btn>
        </v-toolbar-items>
    </v-app-bar>
</template>

<script lang="ts">

import { createContactUsUrl, createLoginUrl, createHomePageUrl } from '../ts/url'

export default {
    props: [
        'companyName'
    ],
    computed:  {
        contactUsUrl: createContactUsUrl,
        homePageUrl: createHomePageUrl,
        loginPageUrl: createLoginUrl
    }
}

</script>

<style scoped>

a {
    text-decoration: none;
    color: black !important;
}

</style>
