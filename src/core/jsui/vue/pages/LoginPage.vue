<template>
    <v-app>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <h1>Login Here</h1>
        </v-content>
    </v-app>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'

export default {
    props: [
        'companyName'
    ],
    components: {
        LandingPageAppBar,
    }
}

</script>
