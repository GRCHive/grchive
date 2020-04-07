<template>
    <section>
        <landing-page-app-bar
            :company-name="companyName"
        >
        </landing-page-app-bar>

        <v-content>
            <div id="mainContent">
                <landing-carousel
                    :company-name="companyName"
                    :imageUrl="bannerImageUrl"
                    tag-line="Governance for the digital age."
                    description="The governance management solution built with big data, cloud vendors, and automation in mind to streamline your GRC process."
                >
                </landing-carousel>

                <v-container class="pa-0" fluid>
                    <feature-highlight
                        headline="Decrease audit lead time"
                        description="Transform narratives into interactive process flows that link directly to controls, risks, documentation, and requests. GRCHive is a centralized solution that seamlessly links all aspects of the audit without you ever having to leave the app. The entire life cycle of the audit is encapsulated to be as efficient as possible to reduce the amount of time and money spent on audits."
                        ctaText="Learn More"
                        :ctaUrl="`${learnMoreUrl}#flow`"
                    >
                        <template v-slot:image>
                            <v-img
                                src="/static/assets/icons/employee.svg"
                                width="250px"
                                max-width="250px"
                            >
                            </v-img>
                        </template>
                    </feature-highlight>

                    <v-divider class="feature-divider"></v-divider>

                    <feature-highlight
                        headline="Maximize audit communication efficiency"
                        description="GRCHive tracks audit document requirements so you don’t have to. Automated notifications and reminders goes out to documentation request owners allowing all key stakeholders to stay on the same page. GRCHive is a central repository that becomes the single source of truth minimizing potential audit status update gaps and ensures key stackholders receive the information they need clearly and succinctly."
                        ctaText="Learn More"
                        :ctaUrl="`${learnMoreUrl}#repo`"
                        flip
                    >
                        <template v-slot:image>
                            <v-img
                                src="/static/assets/icons/negotiating.svg"
                                width="250px"
                                max-width="250px"
                            >
                            </v-img>
                        </template>
                    </feature-highlight>

                    <v-divider class="feature-divider"></v-divider>

                    <feature-highlight
                        headline="Automate operational confidence"
                        description="Allow users to automate testing by creating custom test scripts to run periodically (coming soon). Additionally, integration with databases allows for users to directly execute read only queries to generate data and to facilitate the IPE evaluation process. GRCHive provides you the confidence your organization needs that the data generated and test results produced are complete and accurate."
                        ctaText="Learn More"
                        :ctaUrl="`${learnMoreUrl}#db`"
                    >
                        <template v-slot:image>
                            <v-img
                                src="/static/assets/icons/automate.svg"
                                width="250px"
                                max-width="250px"
                            >
                            </v-img>
                        </template>
                    </feature-highlight>

                    <v-divider class="feature-divider"></v-divider>

                    <feature-highlight
                        headline="Secure audit handoffs"
                        description="Tighten data security by allowing only users who need access to sensitive information with GRCHive’s customizable user roles. GRCHive also features an audit trail to monitor the data request and collection process for key stakeholders to review. Secure your audit data by knowing who, what, when, where and how for all aspects of your audit within GRCHive."
                        ctaText="Learn More"
                        :ctaUrl="`${learnMoreUrl}#security`"
                        flip
                    >
                        <template v-slot:image>
                            <v-img
                                src="/static/assets/icons/security.svg"
                                width="250px"
                                max-width="250px"
                            >
                            </v-img>
                        </template>
                    </feature-highlight>

                    <v-divider class="feature-divider"></v-divider>

                    <standalone-cta class="max-height" id="bottomCta">
                    </standalone-cta>
                </v-container>
            </div>

            <landing-footer
                :flaticon-authors="['xnimrodx']"
            ></landing-footer>
        </v-content>
    </section>
</template>

<script lang="ts">

import LandingPageAppBar from '../components/LandingPageAppBar.vue'
import LandingCarousel from '../components/info/LandingCarousel.vue'
import FeatureHighlight from '../components/info/FeatureHighlight.vue'
import LandingFooter from '../components/info/LandingFooter.vue'
import StandaloneCta from '../components/info/StandaloneCta.vue'
import { PageParamsStore } from '../../ts/pageParams'
import { learnMoreUrl, createAssetUrl } from '../../ts/url'

export default {
    components: {
        LandingPageAppBar,
        LandingCarousel,
        FeatureHighlight,
        LandingFooter,
        StandaloneCta,
    },
    data: () => ({
        learnMoreUrl,
    }),
    computed: {
        companyName() : string {
            return PageParamsStore.state.site!.CompanyName
        },
        bannerImageUrl() : string {
            return createAssetUrl('working.jpg')
        },
    }
}

</script>
