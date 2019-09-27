<template>
    <v-app>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <v-container fluid class="pa-0 max-height">
                <v-row 
                    align="center"
                    justify="center"
                    class="max-height"
                >
                        <v-card width="80%">
                            <v-row
                                align="center"
                            >
                                <v-col cols=8>
                                    <v-card-title>
                                        Thank you for your interest!
                                    </v-card-title>
                                    <v-card-text id="message">
                                        Our engineers are hard at work getting {{ companyName }} ready for you.
                                        Leave us your name and email address, and we will contact you with more information on how we can help you simplify your external and internal auditing processes.
                                    </v-card-text>
                                    <v-col cols=12>
                                        <collect-email :company-name="companyName"></collect-email>
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
    </v-app>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import CollectEmail from '../components/CollectEmail.vue'
import { createAssetUrl} from '../../ts/url'

export default {
    props: [
        'companyName'
    ],
    components: {
        LandingPageAppBar,
        CollectEmail
    },
    computed: {
        imageUrl: () => createAssetUrl("construct.jpg")
    }
}

</script>

<style scoped>

#message {
    word-break: keep-all;
}

</style>
