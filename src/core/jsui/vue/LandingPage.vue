<template>
    <v-app>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <landing-carousel
                :company-name="companyName"
            >
            </landing-carousel>
        </v-content>
    </v-app>
</template>

<script lang="ts">

import LandingPageAppBar from './LandingPageAppBar.vue'
import LandingCarousel from './LandingCarousel.vue'

export default {
    props: [
        'companyName'
    ],
    components: {
        LandingPageAppBar,
        LandingCarousel
    }
}

</script>
