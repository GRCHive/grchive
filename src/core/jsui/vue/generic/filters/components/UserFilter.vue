<template>
    <div class="d-flex">
        <span class="font-weight-bold label-helper mr-4">{{ label }}</span>
        <multi-user-search-form-component
            label="Users"
            :value="userArray"
            @input="onChange"
            :select-noone="selectNoone"
            hide-details
            dense
        >
        </multi-user-search-form-component>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import { UserFilterData } from '../../../../ts/filters'
import MultiUserSearchFormComponent from '../../MultiUserSearchFormComponent.vue'

@Component({
    components: {
        MultiUserSearchFormComponent
    }
})
export default class UserFilter extends Vue {
    @Prop()
    label! : string

    @Prop()
    value!: UserFilterData

    @Prop()
    selectNoone!: boolean

    // TODO: Sync from user ids to userArray. This case isn't currently
    // being used so we can just not implemenet for now.
    userArray : (User | null)[] = []

    onChange( v : (User | null)[]) {
        this.userArray = v
        this.value.UserIds = v.map((ele : User | null) => {
            if (!!ele) {
                return ele.Id
            } else {
                return null
            }
        })
        this.$emit('input', this.value)
    }
}

</script>

<style scoped>

.label-helper { 
    margin-top: 10px;
}

</style>
