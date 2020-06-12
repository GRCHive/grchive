<template>
    <div>
        <v-content>
            <v-container class="pa-0" fluid>
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
                            text="Welcome Back"
                        ></hero-image>

                        <v-row justify="center">
                            <v-card width="30%">
                                <v-card-title>
                                    Login
                                </v-card-title>
                                <login-form></login-form>
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
import HeroImage from '../components/HeroImage.vue'
import LoginForm from '../components/LoginForm.vue'
import { createAssetUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'

export default {
    components: {
        LandingPageAppBar,
        HeroImage,
        LoginForm,
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
