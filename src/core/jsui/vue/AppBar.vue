<template>
    <v-app-bar app>
        <v-toolbar-title>{{ companyName }}</v-toolbar-title>
    </v-app-bar>
</template>

<script lang="ts">

export default {
    props: [
        'companyName'
    ],
}

</script>
