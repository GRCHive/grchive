<template>
    <section>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <hero-image
                :src="bannerImageUrl"
                :max-height=300
                position="center center"
                text="Let's Talk"
            ></hero-image>

            <v-container fluid>
                <v-row justify="center">
                    <v-card width="30%" class="ma-4">
                        <v-card-title>
                            Email Us
                        </v-card-title>
                        <v-card-text class="long-text">
                            <div class="long-text">
                                Send us an email and we will respond as soon as we can!
                            </div>

                            <div class="text--primary headline my-2">
                                Interested in our product?
                                <a :href="salesEmail.mailto">{{ salesEmail.email }}</a>
                            </div>

                            <div class="text--primary headline my-2">
                                Interested in working with us?
                                <a :href="careerEmail.mailto">{{ careerEmail.email }}</a>
                            </div>

                            <div class="text--primary headline my-2">
                                Need technical or audit support?
                                <a :href="supportEmail.mailto">{{ supportEmail.email }}</a>
                            </div>

                            <div class="text--primary headline my-2">
                                All other inquiries:
                                <a :href="contactEmail.mailto">{{ contactEmail.email }}</a>
                            </div>
                        </v-card-text>
                    </v-card>

                    <!-- Hide for now since I don't want to deal with doing the whole SMTP stuff.
                    <v-card width="30%" class="ma-4">
                        <v-card-title>
                            Send us a Message
                        </v-card-title>
                        <v-card-text class="long-text">
                            Can't find an email address that describes your concern?
                            Don't want to open up your email client?
                            No problem!
                            Send us a message using the form below and we will take care of it.
                        </v-card-text>

                        <div class="ma-3">
                            <message-us></message-us>
                        </div>
                    </v-card>
                    -->
                </v-row>
            </v-container>

        </v-content>
    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import HeroImage from '../components/HeroImage.vue'
import MessageUs from '../components/MessageUs.vue'
import { createAssetUrl, createMailtoUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import Vue from 'vue'

export default Vue.extend({
    components: {
        LandingPageAppBar,
        HeroImage,
        MessageUs
    },
    computed: {
        companyName() : string {
            return PageParamsStore.state.site!.CompanyName
        },
        domain() : string {
            return PageParamsStore.state.site!.Domain
        },
        bannerImageUrl() : string {
            return createAssetUrl('generic-banner2.jpg')
        },
        salesEmail() : Object {
            return createMailtoUrl('sales', this.domain)
        },
        careerEmail() : Object {
            return createMailtoUrl('careers', this.domain)
        },
        supportEmail() : Object {
            return createMailtoUrl('support', this.domain)
        },
        contactEmail() : Object {
            return createMailtoUrl('contact', this.domain)
        }
    }
})

</script>
