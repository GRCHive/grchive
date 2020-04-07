<template>
    <v-footer
        padless
    >
        <div class="max-width" id="footerContent">
            <v-row class="caption" align="center" justify="center">
                <v-col cols="3">
                    <p class="info-text">
                        Icons made by 
                        <a v-for="(author, idx) in flaticonAuthors"
                           :href="`https://www.flaticon.com/authors/${author}`" :title="author"
                        >
                            {{ author }}
                        </a>
                        from
                        <a href="https://www.flaticon.com/" title="Flaticon"> www.flaticon.com</a>
                    </p>
                </v-col>
                <v-col cols="7">
                    <v-row justify="end">
                        <p class="info-text">
                            Copyright &copy; 2020 GRCHive, Inc.
                        </p>
                    </v-row>
                </v-col>
            </v-row>
        </div>
    </v-footer>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        flaticonAuthors : Array
    }
})

@Component
export default class LandingFooter extends Props {

}

</script>

<style scoped>

#footerContent {
    background-color: #263238;
    color: white;
}

.info-text {
    color: #BDBDBD;
    margin-bottom: 0px !important;
}

</style>
