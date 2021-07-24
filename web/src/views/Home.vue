<template lang="pug">
van-loading(v-if="loading" style="text-align:center;margin-top:10px;") Loading...
van-cell(v-else v-for="(t, i) in todos" :key="t.id" :title="t.text")
  template(#right-icon)
    van-switch(v-model="doneSwitchState[t.id]" size="24")
</template>

<script>
import { ref } from "vue";
import { useQuery, useResult } from "@vue/apollo-composable";
import gql from "graphql-tag";

export default {
  name: "Home",
  setup() {
    const { result, loading } = useQuery(gql`
      query listTodos {
        todos {
          id
          text
          done
        }
      }
    `);
    const todos = useResult(result, [], (data) => data.todos);
    const doneSwitchState = ref({});
    return { loading, todos, doneSwitchState };
  },
};
</script>
