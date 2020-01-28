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
import { PageParamsStore } from '../../ts/pageParams'

export default {
    components: {
        LandingPageAppBar,
        LandingCarousel
    },
    computed: {
        companyName() : string {
            return PageParamsStore.state.site!.CompanyName
        }
    }
}

</script>
