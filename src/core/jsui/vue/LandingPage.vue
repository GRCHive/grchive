<template>
    <v-app>
        <landing-page-app-bar
            :company-name="companyName">
        </landing-page-app-bar>
    </v-app>
</template>

<script lang="ts">

import HomePageAppBar from './HomePageAppBar.vue'

export default {
    props: [
        'companyName'
    ],
    components: {
        HomePageAppBar
    }
}

</script>
