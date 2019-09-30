<template>
    <v-app>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <hero-image
                :src="bannerImageUrl"
                max-height=300
                position="center center"
                text="Welcome Back"
            ></hero-image>

            <v-container fluid>
                <v-row justify="center">
                    <v-card width="30%">
                        <v-card-title>
                            Login
                        </v-card-title>
                        <login-form :csrf="csrf"></login-form>
                    </v-card>
                </v-row>
            </v-container>

        </v-content>

    </v-app>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import HeroImage from '../components/HeroImage.vue'
import LoginForm from '../components/LoginForm.vue'
import { createAssetUrl } from '../../ts/url'

export default {
    props: [
        'companyName',
        'csrf'
    ],
    components: {
        LandingPageAppBar,
        HeroImage,
        LoginForm,
    },
    computed: {
        bannerImageUrl() : string {
            return createAssetUrl('generic-banner3.jpg')
        },
    }
}

</script>
