<template>
    <div>
        <v-content>
            <v-container fluid class="pa-0">
                <v-row justify="center">
                    <v-col cols="10" id="mainContent" class="pa-0">
                        <landing-page-app-bar
                            :company-name="companyName"
                        >
                        </landing-page-app-bar>

                        <hero-image
                            :src="bannerImageUrl"
                            :max-height=300
                            position="center center"
                            text="Get Started"
                        ></hero-image>

                        <v-row
                            justify="center"
                        >
                            <v-card width="80%" id="card-content">
                                <v-row>
                                    <v-col cols=8>
                                        <v-card-title>
                                            Request a live demo.
                                        </v-card-title>
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
                    </v-col>
                </v-row>
            </v-container>
        </v-content>
    </div>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import CollectEmail from '../components/CollectEmail.vue'
import HeroImage from '../components/HeroImage.vue'
import { createAssetUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

export default {
    components: {
        LandingPageAppBar,
        CollectEmail,
        HeroImage
    },
    computed: {
        imageUrl: () => createAssetUrl("construct.jpg"),
        bannerImageUrl: () => createAssetUrl("generic-banner1.jpg"),
        companyName() : string {
            return PageParamsStore.state.site!.CompanyName
        }
    }
}

</script>
