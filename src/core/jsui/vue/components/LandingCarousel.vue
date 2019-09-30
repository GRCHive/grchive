<template>
    <section>
        <v-img
            :src="imageUrl"
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
                                Streamline your audit process with {{ companyName }}.
                            </p>
                        </v-row>

                        <v-row
                            justify="center"
                            class="ma-2"
                        >
                            <v-btn 
                                class="above-overlay ma-1"
                                color="primary"
                                :href="getStartedUrl"
                            >
                                Get Started
                            </v-btn>

                            <v-btn 
                                class="above-overlay ma-1"
                                color="accent"
                                :href="learnMoreUrl"
                            >
                                Learn More
                            </v-btn>

                        </v-row>

                    </v-col>
                </v-row>
            </v-container>
        </v-img>
    </section>
</template>

<script lang="ts">

import { getStartedUrl, learnMoreUrl, createAssetUrl} from '../../ts/url'

export default {
    props: [
        'companyName'
    ],
    data: () => ({
        getStartedUrl,
        learnMoreUrl
    }),
    computed:  {
        imageUrl: () => createAssetUrl("working.jpg"),
    }
}

</script>
