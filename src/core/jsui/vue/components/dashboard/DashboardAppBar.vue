<template>
    <v-app-bar app dense clipped-left clipped-right :extension-height="extensionHeight">
        <v-app-bar-nav-icon @click.stop="clickNav"></v-app-bar-nav-icon>
        <v-toolbar-title color="primary">
            <a :href="orgUrl">{{ orgName }}</a>
        </v-toolbar-title>
        <div class="flex-grow-1"></div>

        <v-toolbar-items>
            <v-btn text
                   color="primary"
                   :href="feedbackUrl.mailto"
            >
                Feedback
                <v-icon color="primary" small>mdi-email</v-icon>
            </v-btn>
            <v-btn text
                   color="primary"
                   :href="supportUrl.mailto"
            >
                Support
                <v-icon color="primary" small>mdi-email</v-icon>
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

        <template v-slot:extension>
            <verify-email-banner @toggle-banner="recomputeExtensionHeight"></verify-email-banner>
        </template>
    </v-app-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import LocalStorage from '../../../ts/localSettings'
import {createLogoutUrl, createMyAccountUrl, createMailtoUrl } from '../../../ts/url'
import { getCurrentCSRF } from '../../../ts/csrf'
import { PageParamsStore } from '../../../ts/pageParams'
import VerifyEmailBanner from '../VerifyEmailBanner.vue'

export default Vue.extend({
    components: {
        VerifyEmailBanner
    },
    data: function() {
        return {
            logoutUrl : createLogoutUrl(getCurrentCSRF()),
            extensionHeight: 40,
        }
    },
    computed: {
        orgUrl() : string  {
            return PageParamsStore.state.organization!.Url
        },
        orgName() : string  {
            return PageParamsStore.state.organization!.Name
        },
        userFirstName() : string {
            return PageParamsStore.state.user!.FirstName
        },
        userLastName() : string {
            return PageParamsStore.state.user!.LastName
        },
        fullName() : string {
            return this.userFirstName + " " + this.userLastName
        },
        myAccountUrl() : string {
            return createMyAccountUrl(PageParamsStore.state.user!.Id)
        },
        supportUrl() : Object {
            return createMailtoUrl("support", PageParamsStore.state.site!.Domain)
        },
        feedbackUrl() : Object {
            return createMailtoUrl("feedback", PageParamsStore.state.site!.Domain)
        },
    },
    methods: {
        clickNav() {
            LocalStorage.commit('setMiniNavBar', !LocalStorage.state.miniNavBar)
        },
        recomputeExtensionHeight(isBannerShown: boolean) {
            if (isBannerShown) {
                this.extensionHeight = 40
            } else {
                this.extensionHeight = 0
            }
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

>>>.v-toolbar__extension {
    padding: 0;
    margin: 0;
}

</style>
