<template>
    <v-img
        :src="src"
        :max-height="maxHeight"
        :position="position"
    >
        <v-overlay
            absolute
            value=true
        >
        </v-overlay>

        <v-row
            align="center"
            justify="center"
            class="max-height"
        >
            <p 
                class="white--text display-2 font-weight-bold above-overlay"
                id="get-started-label"
            >
                {{ text }}
            </p>
        </v-row>
    </v-img>
</template>

<script lang="ts">

export default {
    props: {
        'position': String,
        'maxHeight': Number,
        'src': String,
        'text': String
    }
}

</script>
