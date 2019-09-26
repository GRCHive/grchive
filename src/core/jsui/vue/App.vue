<template>
    <v-app>
        <AppBar
            :company-name="companyName">
        </AppBar>

    </v-app>
</template>

<script lang="ts">

import AppBar from './AppBar'

export default {
    props: [
        'companyName'
    ],
    components: {
        AppBar
    }
}

</script>
