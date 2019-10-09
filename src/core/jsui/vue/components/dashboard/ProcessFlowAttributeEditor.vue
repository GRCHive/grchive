<template>
    <v-navigation-drawer absolute right :style="clipStyle" ref="attrNavDrawer" :value="showHide">
        <h1>BOOP BOOP I"M A ROBOT</h1>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 

export default Vue.extend({
    props: {
        customClipHeight : Number,
        showHide : Boolean
    },
    computed: {
        clipStyle() {
            return {
                //"transform": "translateX(0%)",
                //"width": "256px",
                "height":  "100vh !important",
                //@ts-ignore
                "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
                //@ts-ignore
                "top" : this.customClipHeight.toString() + "px"
            }
        },
    }
})

</script>
