<template>
    <section>
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
    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import LandingCarousel from '../components/LandingCarousel.vue'

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
