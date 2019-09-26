<template>
    <v-app-bar app>
        <v-toolbar-title color="primary">{{ companyName }}</v-toolbar-title>
        <div class="flex-grow-1"></div>

        <v-toolbar-items>
            <v-btn text color="primary">Contact Us</v-btn>
            <v-btn text color="primary">Login</v-btn>
        </v-toolbar-items>
    </v-app-bar>
</template>

<script lang="ts">

export default {
    props: [
        'companyName'
    ]
}

</script>
