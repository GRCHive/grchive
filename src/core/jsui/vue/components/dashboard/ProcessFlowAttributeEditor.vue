<template>
    <v-navigation-drawer absolute permanent right :style="clipStyle" ref="attrNavDrawer">
        <h1>BOOP BOOP I"M A ROBOT</h1>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 

export default Vue.extend({
    props: {
        customClipHeight : Number
    },
    computed: {
        clipStyle() {
            return {
                "transform": "translateX(0%)",
                "width": "256px",
                "height":  "100vh !important",
                //@ts-ignore
                "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
                //@ts-ignore
                "top" : this.customClipHeight.toString() + "px"
            }
        },
    },
    methods: {
        refreshStyle() {
            //@ts-ignore
            Object.keys(this.clipStyle).forEach((item) => {
                //@ts-ignore
                this.$refs.attrNavDrawer.$el.style[item] = this.clipStyle[item]
            })
        }
    },
    mounted() {
        //@ts-ignore
        this.refreshStyle()
    }
})

</script>
