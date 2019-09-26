<template>
    <section>
        <v-img
            src="static/assets/working.jpg"
            max-height=500
            position="top center"
        >
            <v-overlay
                absolute
                value=true
            >
            </v-overlay>

            <v-container fluid class="max-height">
                <v-row
                    align="center"
                    class="max-height"
                >
                    <v-col
                        cols=12
                    >
                        <v-row
                            justify="center"
                            class="ma-2"
                        >
                            <p
                                class="white--text display-2 font-weight-bold above-overlay"
                            >
                                Transform your audit process with {{ companyName }}.
                            </p>
                        </v-row>

                        <v-row
                            justify="center"
                            class="ma-2"
                        >
                            <v-btn class="above-overlay" color="primary">Get Started</v-btn>
                        </v-row>

                    </v-col>
                </v-row>
            </v-container>
        </v-img>
    </section>
</template>

<script lang="ts">

export default {
    props: [
        'companyName'
    ]
}

</script>

<style scoped>

.above-overlay {
    z-index: 7;
}

</style>
