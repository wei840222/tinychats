<template lang="pug">
div {{ accessToken }}
router-view
</template>

<script>
import { ref, provide, onMounted } from "vue";
import {
  ApolloClient,
  createHttpLink,
  InMemoryCache,
} from "@apollo/client/core";
import { createPersistedQueryLink } from "@apollo/client/link/persisted-queries";
import { sha256 } from "crypto-hash";
import { DefaultApolloClient } from "@vue/apollo-composable";
import liff from "@line/liff";

export default {
  setup() {
    const apolloClient = new ApolloClient({
      link: createPersistedQueryLink({
        useGETForHashedQueries: true,
        sha256,
      }).concat(
        createHttpLink({
          uri: "https://wei840222-todo.herokuapp.com/graphql",
        })
      ),
      cache: new InMemoryCache(),
    });
    provide(DefaultApolloClient, apolloClient);

    const accessToken = ref("");

    onMounted(async () => {
      await liff.init({ liffId: "1656247924-eX5ZOvN0" });
      if (!liff.isLoggedIn()) {
        liff.login();
        return;
      }
      accessToken.value = liff.getAccessToken();
    });

    return { accessToken };
  },
};
</script>

<style lang="sass">
#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  color: #2c3e50
</style>
