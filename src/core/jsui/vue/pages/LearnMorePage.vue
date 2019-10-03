<template>
    <section>
        <landing-page-app-bar
            :company-name="this.$root.companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <v-img
                :src="bannerImageUrl"
                max-height=600
                position="center center"
            >
                <v-overlay
                    absolute
                    value=true
                    opacity=0.7
                >
                </v-overlay>

                <v-row
                    align="center"
                    class="max-height"
                >
                    <v-col cols=12>
                        <v-row justify="center" class="ma-4">
                            <div class="white--text display-2 font-weight-bold above-overlay">
                                What is {{ this.$root.companyName }}?
                            </div>
                        </v-row>

                        <v-row justify="center" class="ma-4">
                            <div class="white--text subtitle-2 above-overlay" id="subtitle">
                                {{ this.$root.companyName }} is an automated audit control management solution that digitizes and streamlines your audit process.
                                It organizes and centralizes information about process flows and audit controls so that everyone, auditors and business managers, can obtain the information they need, faster.
                            </div>
                        </v-row>

                        <v-row justify="center" class="ma-4">
                            <v-btn 
                                class="above-overlay ma-1"
                                color="primary"
                                :href="getStartedUrl"
                            >
                                Get Started
                            </v-btn>

                            <v-btn 
                                class="above-overlay ma-1"
                                color="accent"
                                :href="contactUsUrl"
                            >
                                Contact Us
                            </v-btn>
                        </v-row>
                    </v-col>
                </v-row>
            </v-img>

            <v-container fluid>

            </v-container>

        </v-content>
    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import { createAssetUrl, getStartedUrl, contactUsUrl } from '../../ts/url'
import Vue from 'vue'

export default Vue.extend({
    components: {
        LandingPageAppBar,
    },
    computed: {
        bannerImageUrl() : string {
            return createAssetUrl('generic-banner4.jpg')
        },
    },
    data : () => ({
        getStartedUrl,
        contactUsUrl
    })
})

</script>

<style scoped>

#subtitle {
    width: 40%;
    text-align: center;
}

</style>
