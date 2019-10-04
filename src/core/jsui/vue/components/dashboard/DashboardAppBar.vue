<template>
    <v-app-bar app dense>
        <v-toolbar-title color="primary">
            <a :href="this.$root.orgUrl">{{ this.$root.orgName }}</a>
        </v-toolbar-title>
        <div class="flex-grow-1"></div>

        <v-toolbar-items>
            <v-btn text
                   color="primary"
                   :href="supportUrl.mailto"
            >
                Support
            </v-btn>
            <v-menu offset-y>

                <template v-slot:activator="{ on }">
                    <v-btn text
                           color="primary"
                           v-on="on"
                    >
                        {{ fullName }}
                        <v-icon color="primary">mdi-menu-down</v-icon>
                    </v-btn>
                </template>
                <v-list dense>
                    <v-list-item dense :href="myAccountUrl">
                        <v-list-item-title>My Account</v-list-item-title>
                    </v-list-item>
                    <v-list-item dense :href="logoutUrl">
                        <v-list-item-title>Logout</v-list-item-title>
                    </v-list-item>
                </v-list>

            </v-menu>
        </v-toolbar-items>
    </v-app-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import {createLogoutUrl, createMyAccountUrl, createMailtoUrl } from '../../../ts/url'

export default Vue.extend({
    data: function() {
        return {
            //@ts-ignore
            fullName: this.$root.userFullName,
            //@ts-ignore
            logoutUrl : createLogoutUrl(this.$root.csrf),
            //@ts-ignore
            myAccountUrl: createMyAccountUrl(this.$root.userEmail),
            //@ts-ignore
            supportUrl: createMailtoUrl("support", this.$root.domain),
        }
    }
})

</script>

<style scoped>

a {
    text-decoration: none;
    color: black !important;
}

.v-menu__content {
    border-radius: 0px !important;
}

</style>
