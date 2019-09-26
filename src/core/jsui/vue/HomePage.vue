<template>
    <v-app>
        <HomePageAppBar
            :company-name="companyName">
        </HomePageAppBar>
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
