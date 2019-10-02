<template>
    <section>
        <landing-page-app-bar
            :company-name="this.$root.companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <landing-carousel
                :company-name="this.$root.companyName"
            >
            </landing-carousel>
        </v-content>
    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import LandingCarousel from '../components/LandingCarousel.vue'

export default {
    components: {
        LandingPageAppBar,
        LandingCarousel
    }
}

</script>
