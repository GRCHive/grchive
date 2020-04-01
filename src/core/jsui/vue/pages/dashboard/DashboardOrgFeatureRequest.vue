<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content class="max-height">
            <v-row class="max-height" align="center">
                <v-col>
                    <v-row justify="center">
                        <p class="display-1">This feature requires an admin to manually enable the <b>{{ featureName }}</b> feature.</p>
                    </v-row>
                    <v-row justify="center" v-if="!localPending">
                        <v-btn
                            color="success"
                            @click="enable"
                        >
                            Enable
                        </v-btn>
                    </v-row>

                    <v-row justify="center" v-else>
                        <v-col cols="12">
                            <p class="subtitle-1 text-center">Your request is in progress. Check back soon!</p>
                            <p class="subtitle-1 text-center">If your request is not processed shortly, please refresh and contact us.</p>
                        </v-col>

                        <v-col cols="6">
                            <v-progress-linear
                                color="primary"
                                indeterminate
                                rounded
                            >
                            </v-progress-linear>
                        </v-col>
                    </v-row>
                </v-col>
            </v-row>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import { enableFeature } from '../../../ts/api/apiFeatures'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        featureName: {
            type: String
        },
        featureId: {
            type: Number
        },
        pending: {
            type: Boolean
        }
    }
})

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
    }
})
export default class DashboardOrgFeatureRequest extends Props {
    localPending : boolean = false

    mounted() {
        this.localPending = this.pending
    }

    enable() {
        this.localPending = true
        enableFeature({
            orgId: PageParamsStore.state.organization!.Id,
            featureId: this.featureId,
        }).then(() => {
            window.location.reload()
        }).catch(() => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }
}
</script>
