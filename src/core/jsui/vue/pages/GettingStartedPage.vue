<template>
    <section>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content class="hide-overflow">
            <hero-image
                :src="bannerImageUrl"
                max-height=300
                position="center center"
                text="Get Started"
            ></hero-image>

            <v-container fluid class="pa-0" id="content">
                <v-row 
                    align="center"
                    justify="center"
                    class="max-height"
                >
                    <v-card width="80%" id="card-content" elevation=12 raised>
                        <v-row
                            align="center"
                        >
                            <v-col cols=8>
                                <v-card-title>
                                    Thank you for your interest!
                                </v-card-title>
                                <v-card-text class="long-text">
                                    Our engineers are hard at work getting {{ companyName }} ready for you.
                                    Leave us your name and email address, and we will contact you with more information on how we can help you streamline your external and internal auditing processes.
                                </v-card-text>
                                <v-col cols=12>
                                    <collect-email 
                                        :company-name="companyName"
                                    ></collect-email>
                                </v-col>
                            </v-col>
                            <v-col
                                cols=4
                                class="pa-0 hide-overflow"
                            >
                                <v-img 
                                    :src="imageUrl"
                                    height="100%"
                                    class="pa-0"
                                ></v-img>
                            </v-col>
                        </v-row>
                    </v-card>
                </v-row>
            </v-container>
        </v-content>
    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import CollectEmail from '../components/CollectEmail.vue'
import HeroImage from '../components/HeroImage.vue'
import { createAssetUrl } from '../../ts/url'

export default {
    props: [
        'companyName'
    ],
    components: {
        LandingPageAppBar,
        CollectEmail,
        HeroImage
    },
    computed: {
        imageUrl: () => createAssetUrl("construct.jpg"),
        bannerImageUrl: () => createAssetUrl("generic-banner1.jpg")
    }
}

</script>

<style scoped>

#content {
    height: calc(100% - 300px);
}

#card-content {
    transform: translateY(-80px);
}

</style>
