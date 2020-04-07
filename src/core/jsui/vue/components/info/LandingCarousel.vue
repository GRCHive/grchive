<template>
    <div>
        <v-img
            :src="imageUrl"
            max-height=500
            :position="position"
        >
            <v-overlay
                absolute
                :value="true"
                :opacity="0.70"
            >
            </v-overlay>

            <v-container fluid class="max-height">
                <v-row
                    align="center"
                    justify="center"
                    class="max-height"
                >
                    <v-col
                        cols=7
                    >
                        <v-row
                            justify="start"
                            class="ma-2"
                        >
                            <p
                                class="white--text above-overlay"
                            >
                                <span class="font-weight-bold display-2">
                                    {{ companyName }}.
                                </span>

                                <span class="display-1">
                                    {{ tagLine }}
                                </span>
                            </p>

                            <p
                                class="white--text above-overlay headline"
                            >
                                {{ description }}
                            </p>
                        </v-row>

                        <v-row
                            justify="start"
                            class="ma-2"
                        >
                            <v-btn 
                                class="above-overlay"
                                color="primary"
                                :href="getStartedUrl"
                                x-large
                            >
                                Request Demo
                            </v-btn>

                            <v-btn 
                                class="above-overlay ml-4"
                                color="accent"
                                :href="learnMoreUrl"
                                x-large
                                v-if="!hideLearnMore"
                            >
                                Learn More
                            </v-btn>

                        </v-row>

                    </v-col>
                </v-row>
            </v-container>
        </v-img>
    </div>
</template>

<script lang="ts">

import { getStartedUrl, learnMoreUrl } from '../../../ts/url'

export default {
    data: () => ({
        getStartedUrl,
        learnMoreUrl
    }),
    props: {
        companyName: String,
        tagLine: String,
        description : String,
        imageUrl: String,
        hideLearnMore : {
            type: Boolean,
            default: false,
        },
        position: {
            type: String,
            default: "top center",
        }
    },
}

</script>
