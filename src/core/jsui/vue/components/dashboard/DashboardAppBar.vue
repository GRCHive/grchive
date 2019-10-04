<template>
    <v-app-bar app dense>
        <v-toolbar-title color="primary">
            <a :href="this.$root.orgUrl">{{ this.$root.orgName }}</a>
        </v-toolbar-title>
        <div class="flex-grow-1"></div>

        <v-toolbar-items>

            <v-menu offset-y>
                <template v-slot:activator="{ on }">
                    <v-btn text
                           color="primary"
                           v-on="on"
                    >
                        {{ fullName }}
                        <v-icon color="primary">mdi-menu-down</v-icon>
                    </v-btn>
                </template>
                <v-list>
                    <v-list-item>
                        <v-list-item-title>Profile</v-list-item-title>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title>Settings</v-list-item-title>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title>Logout</v-list-item-title>
                    </v-list-item>
                </v-list>

            </v-menu>
        </v-toolbar-items>
    </v-app-bar>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
    data: function() {
        return {
            //@ts-ignore
            fullName: this.$root.userFullName,
        }
    }
})

</script>

<style scoped>

a {
    text-decoration: none;
    color: black !important;
}

</style>
