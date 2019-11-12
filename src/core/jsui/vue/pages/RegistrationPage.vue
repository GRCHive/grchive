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
                text="Welcome"
            ></hero-image>

            <v-container fluid>
                <v-row justify="center">
                    <v-card width="30%">
                        <v-card-title>
                            Register
                        </v-card-title>
                        <registration-form></registration-form>
                    </v-card>
                </v-row>
            </v-container>
        </v-content>

    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import HeroImage from '../components/HeroImage.vue'
import RegistrationForm from '../components/RegistrationForm.vue'
import { createAssetUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

export default {
    components: {
        LandingPageAppBar,
        HeroImage,
        RegistrationForm
    },
    computed: {
        bannerImageUrl() : string {
            return createAssetUrl('generic-banner3.jpg')
        },
        companyName() : string {
            return PageParamsStore.state.site!.CompanyName
        }
    }
}

</script>
