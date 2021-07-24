<template lang="pug">
div {{ JSON.stringify(result) }}
</template>

<script>
import { useQuery } from "@vue/apollo-composable";
import gql from "graphql-tag";

export default {
  name: "Home",
  setup() {
    const { result } = useQuery(gql`
      query listTodos {
        todos {
          id
          text
          done
        }
      }
    `);
    return {
      result,
    };
  },
};
</script>
